/*
 * Copyright 2014–2018 Kullo GmbH
 *
 * This source code is licensed under the 3-clause BSD license. See LICENSE.txt
 * in the root directory of this source tree for details.
 */
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

var crashIds = make(chan string, 100)

func readAllFromCommand(cmd *exec.Cmd) ([]byte, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return stdoutBytes, nil
}

func pullGit(repository string) error {
	cmd := exec.Command("git", "-C", repository, "pull")
	return cmd.Run()
}

func StartCrashreportProcessorWorker() {
	go func() {
		linux32Symbols := *symbolsDirectory + "/linux32"
		linux64Symbols := *symbolsDirectory + "/linux64"
		osxSymbols := *symbolsDirectory + "/osx"
		windowsSymbols := *symbolsDirectory + "/windows"

		for {
			id := <-crashIds
			baseFilename := *dumpDirectory + "/" + id
			//metadataFilename := baseFilename + ".json"
			minidumpFilename := baseFilename + ".dmp"
			stacktraceFilename := baseFilename + ".trace"

			err := pullGit(*symbolsDirectory)
			if err != nil {
				log.Print("Failed to pull git")
			}

			cmd := exec.Command("minidump_stackwalk", minidumpFilename,
				linux32Symbols, linux64Symbols, osxSymbols, windowsSymbols)
			stdout, err := readAllFromCommand(cmd)
			if err != nil {
				log.Print(fmt.Sprintf("Failed to process minidump %s", minidumpFilename))
				continue
			}

			err = ioutil.WriteFile(stacktraceFilename, stdout, 0600)
			if err != nil {
				log.Print(fmt.Sprintf("Failed to write stacktrace %s", stacktraceFilename))
				continue
			}
		}
	}()
}

func ProcessCrashreport(id string) {
	crashIds <- id
}

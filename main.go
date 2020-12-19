/*
 * Copyright 2014–2018 Kullo GmbH
 *
 * This source code is licensed under the 3-clause BSD license. See LICENSE.txt
 * in the root directory of this source tree for details.
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"runtime/debug"
)

var listen = flag.String("listen", ":8080", "where to listen (address and port)")
var errorLogFile = flag.String("error-logfile", "", "file name of the error log")
var dumpDirectory = flag.String("dump-directory", "/tmp", "location of dumps and corresponding data")
var symbolsDirectory = flag.String("symbols-directory", "/opt/breakpad-symbols", "location of the symbols repository")

func writeServerError(err error, rw http.ResponseWriter) {
	log.Print(err.Error() + "\n" + string(debug.Stack()))
	rw.WriteHeader(http.StatusInternalServerError)
	io.WriteString(rw, "Internal Server Error. We're sorry.")
}

func statusHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, "HTTP 200\n\nUp and running")
}

func uploadHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "Bad request: This endpoint only allows POST requests")
		return
	}

	contentType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "Bad request: Content-Type cannot be parsed")
		return
	}

	if contentType != "multipart/form-data" {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "Bad request: Content-Type must be multipart/form-data")
		return
	}

	mr := multipart.NewReader(req.Body, params["boundary"])
	form, err := mr.ReadForm(10 * 1024 * 1024)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "Bad request: Multi-part message cannot be parsed")
		return
	}

	minidumpHeader := form.File["upload_file_minidump"]
	if minidumpHeader == nil {
		rw.WriteHeader(http.StatusBadRequest)
		io.WriteString(rw, "Bad request: Body must contain upload_file_minidump")
		return
	}

	uploadId, err := writeToDisk(dumpDirectory, minidumpHeader, form.Value)
	if err != nil {
		writeServerError(err, rw)
		return
	}

	ProcessCrashreport(uploadId)

	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, uploadId)
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	flag.Parse()

	if *errorLogFile != "" {
		openedErrorLogFile, err := os.OpenFile(
			*errorLogFile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0640)
		if err != nil {
			log.Fatal("Couldn't open error log: " + err.Error())
		}
		log.SetOutput(openedErrorLogFile)
	}

	StartCrashreportProcessorWorker()

	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/upload", uploadHandler)

	log.Print(fmt.Sprintf("Starting HTTP server on '%s' ...", *listen))
	log.Fatal(http.ListenAndServe(*listen, nil))
}

/*
 * Copyright 2014–2018 Kullo GmbH
 *
 * This source code is licensed under the 3-clause BSD license. See LICENSE.txt
 * in the root directory of this source tree for details.
 */
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
)

type Metadata struct {
	Prod     string
	Ver      string
	Guid     string
	Ptime    string
	Ctime    string
	Email    string
	Comments string
}

func writeToDisk(dumpDirectory *string, fileHeader []*multipart.FileHeader, attribs map[string][]string) (string, error) {
	var meta Metadata
	meta.Prod = getFirstOrEmptyString(attribs["prod"])
	meta.Ver = getFirstOrEmptyString(attribs["ver"])
	meta.Guid = getFirstOrEmptyString(attribs["guid"])
	meta.Ptime = getFirstOrEmptyString(attribs["ptime"])
	meta.Ctime = getFirstOrEmptyString(attribs["ctime"])
	meta.Email = getFirstOrEmptyString(attribs["email"])
	meta.Comments = getFirstOrEmptyString(attribs["comments"])

	jsonMeta, err := json.Marshal(meta)
	if err != nil {
		return "", err
	}
	const NEWLINE = byte(0x0A)
	jsonMeta = append(jsonMeta, NEWLINE)

	basename := getRandomFilename()
	err = ioutil.WriteFile(*dumpDirectory+"/"+basename+".json", jsonMeta, 0600)
	if err != nil {
		return "", err
	}
	dumpFile, err := fileHeader[0].Open()
	if err != nil {
		return "", err
	}
	dump, err := ioutil.ReadAll(dumpFile)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(*dumpDirectory+"/"+basename+".dmp", dump, 0600)
	if err != nil {
		return "", err
	}

	return basename, nil
}

func getFirstOrEmptyString(strings []string) string {
	if len(strings) == 0 {
		return ""
	}
	return strings[0]
}

func getRandomFilename() string {
	rndBytes := make([]byte, 16)
	rand.Read(rndBytes)
	return hex.EncodeToString(rndBytes)
}

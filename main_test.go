/*
 * Copyright 2014–2018 Kullo GmbH
 *
 * This source code is licensed under the 3-clause BSD license. See LICENSE.txt
 * in the root directory of this source tree for details.
 */
package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
)

type ResponseWriterStub struct {
	Status       int
	BytesWritten int
}

func (rws *ResponseWriterStub) Header() http.Header {
	return http.Header{}
}

func (rws *ResponseWriterStub) Write(bytes []byte) (int, error) {
	rws.BytesWritten += len(bytes)
	return len(bytes), nil
}

func (rws *ResponseWriterStub) WriteHeader(status int) {
	rws.Status = status
}

func TestStatus(t *testing.T) {
	var body bytes.Buffer
	req, err := http.NewRequest("GET", "/status", &body)
	if err != nil {
		t.Fatal("http.NewRequest")
	}
	rw := ResponseWriterStub{}

	statusHandler(&rw, req)

	if rw.Status != http.StatusOK {
		t.Error("Status is", rw.Status)
	}
}

func TestUploadHandlerFailsOnGet(t *testing.T) {
	var body bytes.Buffer
	req, err := http.NewRequest("GET", "/upload", &body)
	if err != nil {
		t.Fatal("http.NewRequest")
	}
	rw := ResponseWriterStub{}

	uploadHandler(&rw, req)

	if rw.Status != http.StatusBadRequest {
		t.Error("Status is", rw.Status)
	}
}

func TestUploadHandlerFailsOnMissingContentType(t *testing.T) {
	var body bytes.Buffer
	req, err := http.NewRequest("POST", "/upload", &body)
	if err != nil {
		t.Fatal("http.NewRequest")
	}
	rw := ResponseWriterStub{}

	uploadHandler(&rw, req)

	if rw.Status != http.StatusBadRequest {
		t.Error("Status is", rw.Status)
	}
}

func TestUploadHandlerFailsOnWrongContentType(t *testing.T) {
	var body bytes.Buffer
	req, err := http.NewRequest("POST", "/upload", &body)
	if err != nil {
		t.Fatal("http.NewRequest")
	}
	req.Header.Set("Content-Type", "text/html")
	rw := ResponseWriterStub{}

	uploadHandler(&rw, req)

	if rw.Status != http.StatusBadRequest {
		t.Error("Status is", rw.Status)
	}
}

func TestUploadHandlerFailsOnMissingFile(t *testing.T) {
	var body bytes.Buffer
	mpWriter := multipart.NewWriter(&body)

	req, err := http.NewRequest("POST", "/upload", &body)
	if err != nil {
		t.Fatal("http.NewRequest")
	}
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	rw := ResponseWriterStub{}

	uploadHandler(&rw, req)

	if rw.Status != http.StatusBadRequest {
		t.Error("Status is", rw.Status)
	}
}

func TestUploadHandlerSucceedsOnUpload(t *testing.T) {
	file := bytes.NewBufferString("I am a minidump")

	var body bytes.Buffer
	mpWriter := multipart.NewWriter(&body)
	formFile, err := mpWriter.CreateFormFile("upload_file_minidump", "test.dmp")
	if err != nil {
		t.Fatal("mpWriter.CreateFormFile")
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		t.Fatal("io.Copy")
	}
	mpWriter.Close()

	req, err := http.NewRequest("POST", "/upload", &body)
	if err != nil {
		t.Fatal("http.NewRequest")
	}
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	rw := ResponseWriterStub{}

	uploadHandler(&rw, req)

	if rw.Status != http.StatusOK {
		t.Error("Status is", rw.Status)
	}

	if rw.BytesWritten == 0 {
		t.Error("Received empty response")
	}
}

package services

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createTempFileInput = []struct {
	input   []string
	isError bool
}{
	{[]string{"fine", "txt"}, false},
	{[]string{"create_failed", "/"}, true},
	{[]string{"/+/", "/"}, true},
}

func TestCreateTempFile(t *testing.T) {
	for _, item := range createTempFileInput {
		_, err := createTempFile(item.input[0], item.input[1])
		assert.Equalf(t, item.isError, err != nil, "Expected %v got %v", item.isError, err != nil)
	}
}

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("Error occured while testing: ", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IndexHandler)
	handler.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "Status code expected: %v got %v", http.StatusOK, rr.Code)
}

var modifyHandlerInput = []struct {
	inputURL       string
	expectedStatus int
}{
	{"/modify/dog-testing.png?mode=3&n=1", 302},
	{"/modify/dog-testing.png?mode=3", 200},
	{"/modify/dog-testing.png", 200},
	{"/modify/invalid.png?mode=3", 400},
	{"/modify/dog-testing.png?mode=abc", 400},
	{"/modify/dog-testing.png?mode=3&n=abc", 400},
}

func TestModifyHandler(t *testing.T) {
	for _, item := range modifyHandlerInput {
		req, err := http.NewRequest("GET", item.inputURL, nil)
		if err != nil {
			t.Error("Error occured while testing: ", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ModifyHandler)
		handler.ServeHTTP(rr, req)
		assert.Equalf(t, item.expectedStatus, rr.Code, "Status code expected: %v got %v", item.expectedStatus, rr.Code)
	}
}

var uploadHandlerInput = []struct {
	testCase       string
	paramName      string
	fileName       string
	expectedStatus int
}{
	{"tc1", "image", "./img/dog-testing.png", 302},
	{"tc2", "non-image", "./img/dog-testing.png", 400},
	{"tc3", "image", "./img/dog-testing.png", 500},
	{"tc4", "image", "./img/dog-testing.png", 500},
}

func TestUploadHandler(t *testing.T) {
	tempIoCopyVar := ioCopyVar
	tempCreateTempFileVar := createTempFileVar
	for _, item := range uploadHandlerInput {
		imgFile, err := os.OpenFile(item.fileName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			t.Error("Failed to Open the file")
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(item.paramName, imgFile.Name())
		if err != nil {
			t.Error("Error occured while creating form file: ", err)
		}

		_, err = io.Copy(part, imgFile)
		if err != nil {
			t.Error("Failed to copy:", err)
		}

		err = writer.Close()
		if err != nil {
			t.Error("Failed while closing writer: ", err)
		}

		req, err := http.NewRequest("POST", "/upload", body)
		if err != nil {
			t.Error("Error occured while testing: ", err)
		}

		if item.testCase == "tc3" {
			createTempFileVar = func(name, ext string) (*os.File, error) {
				return nil, errors.New("Temp file failed to create")
			}
		}
		if item.testCase == "tc4" {
			ioCopyVar = func(dst io.Writer, src io.Reader) (written int64, err error) {
				return -1, errors.New("IoCopy failed")
			}
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(UploadHandler)
		handler.ServeHTTP(rr, req)
		assert.Equalf(t, item.expectedStatus, rr.Code, "Status code expected: %v got %v", item.expectedStatus, rr.Code)
		ioCopyVar = tempIoCopyVar
		createTempFileVar = tempCreateTempFileVar
		imgFile.Close()
	}

}

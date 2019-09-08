package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/chroma"
	"github.com/stretchr/testify/assert"
)

func TestPanichandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Error("Error occured while testing: ", err)
	}
	rr := httptest.NewRecorder()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic happpend:", err)
		}
	}()
	handler := http.HandlerFunc(PanicHandler)
	handler.ServeHTTP(rr, req)
	//PanicHandler(rr, req)
	assert.Equalf(t, http.StatusInternalServerError, rr.Code, "After Panic: Status code expected: %v got %v", http.StatusInternalServerError, rr.Code)

}

func TestSourceCodeHandler(t *testing.T) {
	var sourceCodeHandlerInput = []struct {
		testCase       string
		inputURL       string
		expectedStatus int
	}{
		{"successful", "/debug?line=24&path=/usr/local/go/src/runtime/debug/stack.go", 200},
		{"file_unavailable", "/debug?path=/home/rajat/go/src/github.com/gophercises/recover/main1.go", 500},
		{"file_unavailable", "/debug", 500},
		{"successful", "/debug?line=24&path=/home/rajat/go/src/github.com/gophercises/recover/main.go", 200},
		{"io_copy_error", "/debug?line=2&path=/usr/local/go/src/runtime/debug/stack.go", 500},
		{"successful", "/debug?line=2&path=/home/rajat/go/src/github.com/gophercises/recover/main.go", 200},
		{"lexer_tokenisation_error", "/debug?line=3&path=/usr/local/go/src/runtime/debug/stack.go", 500},
	}

	tempIoCopy := IoCopy
	tempLexerTokenise := LexerTokenise

	for _, item := range sourceCodeHandlerInput {
		if item.testCase == "io_copy_error" {
			IoCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
				return -1, errors.New("IoCopy failed")
			}
		}

		if item.testCase == "lexer_tokenisation_error" {
			LexerTokenise = func(options *chroma.TokeniseOptions, text string) (chroma.Iterator, error) {
				return nil, errors.New("Tokenisation failed")
			}
		}

		req, err := http.NewRequest("GET", item.inputURL, nil)
		if err != nil {
			t.Error("Error occured while testing: ", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(SourceCodeHandler)
		handler.ServeHTTP(rr, req)
		//SourceCodeHandler(rr, req)
		assert.Equalf(t, item.expectedStatus, rr.Code, "Status code expected: %v got %v", item.expectedStatus, rr.Code)
		IoCopy = tempIoCopy
		LexerTokenise = tempLexerTokenise
	}
}

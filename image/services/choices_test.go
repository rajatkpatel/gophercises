package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gophercises/image/primitive"
	"github.com/stretchr/testify/assert"
)

var choicesInput = []struct {
	testCase       string
	input          []string
	expectedStatus int
}{
	{"tc1", []string{"../temp/test_input.png", "png"}, 200},
	{"tc2", []string{"../temp/test_input.png", "png"}, 500},
	{"tc3", []string{"../temp/test_input.png", "png"}, 200},
	{"tc4", []string{"../temp/test_input.png", "png"}, 200},
}

func TestRenderShapeChoices(t *testing.T) {
	tempTplExecuteBool := tplExecuteBool
	for _, item := range choicesInput {
		req, err := http.NewRequest("GET", "localhost:5000", nil)
		if err != nil {
			t.Error("Error occured while testing: ", err)
		}
		rr := httptest.NewRecorder()
		imgFile, err := os.OpenFile(item.input[0], os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			t.Error("Failed to Open the file")
		}
		if item.testCase == "tc2" {
			imgFile.Close()
		}
		if item.testCase == "tc4" {
			tplExecuteBool = true
		}
		renderShapeChoices(rr, req, imgFile, item.input[1], primitive.ModeCircle)
		assert.Equalf(t, item.expectedStatus, rr.Code, "Status code expected: %v got %v", item.expectedStatus, rr.Code)
		tplExecuteBool = tempTplExecuteBool
		imgFile.Close()
	}
}

func TestRenderModeChoices(t *testing.T) {
	tempTplExecuteBool := tplExecuteBool
	for _, item := range choicesInput {
		req, err := http.NewRequest("GET", "localhost:5000", nil)
		if err != nil {
			t.Error("Error occured while testing: ", err)
		}
		rr := httptest.NewRecorder()
		imgFile, err := os.OpenFile(item.input[0], os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			t.Error("Failed to Open the file")
		}
		if item.testCase == "tc2" {
			imgFile.Close()
		}
		if item.testCase == "tc4" {
			tplExecuteBool = true
		}
		renderModeChoices(rr, req, imgFile, item.input[1])
		assert.Equalf(t, item.expectedStatus, rr.Code, "Status code expected: %v got %v", item.expectedStatus, rr.Code)
		tplExecuteBool = tempTplExecuteBool
		imgFile.Close()
	}
}

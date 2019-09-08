package services

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/gophercises/image/primitive"
	"github.com/stretchr/testify/assert"
)

var generateImagesInput = []struct {
	testCase string
	input    []string
	isError  bool
}{
	{"tc1", []string{"../temp/test_input.png", "png"}, false},
	{"tc2", []string{"invalid_img_input.png", "png"}, true},
	{"tc3", []string{"invalid_img_input.png", "/"}, true},
	{"tc4", []string{"../temp/test_input.png", "png"}, true},
}

func TestGenerateImages(t *testing.T) {
	option := []generateOptions{
		{N: 10, M: primitive.ModeCircle},
		{N: 10, M: primitive.ModeCombo},
	}

	for _, item := range generateImagesInput {
		imgFile, err := os.OpenFile(item.input[0], os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			t.Error("Failed to Open the file")
		}
		if item.testCase == "tc4" {
			imgFile.Close()
		}
		_, err = generateImages(imgFile, item.input[1], option...)
		assert.Equalf(t, item.isError, err != nil, "Expected %v got %v", item.isError, err != nil)
		imgFile.Close()

	}
}

var imageInput = []struct {
	testCase string
	input    []string
	isError  bool
}{
	{"tc1", []string{"../temp/test_input.png", "png"}, false},
	{"tc2", []string{"../temp/test_input.png", "png"}, true},
	{"tc3", []string{"../temp/test_input.png", "png"}, true},
	{"tc4", []string{"../temp/test_input.png", "png"}, true},
}

func TestGenerateImage(t *testing.T) {
	tempNameVar := TempName
	tempIoCopyVar := ioCopyVar
	for _, item := range imageInput {
		imgFile, err := os.OpenFile(item.input[0], os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			t.Error("Failed to Open the file")
		}
		if item.testCase == "tc2" {
			TempName = "/+/"
		}
		if item.testCase == "tc3" {
			ioCopyVar = func(dst io.Writer, src io.Reader) (written int64, err error) {
				return -1, errors.New("IoCopy failed")
			}
		}
		if item.testCase == "tc4" {
			imgFile.Close()
		}
		_, err = generateImage(imgFile, item.input[1], 5, primitive.ModeCircle)
		assert.Equalf(t, item.isError, err != nil, "Expected %v got %v", item.isError, err != nil)
		TempName = tempNameVar
		ioCopyVar = tempIoCopyVar
		imgFile.Close()
	}
}

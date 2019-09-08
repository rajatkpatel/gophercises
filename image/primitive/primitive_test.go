package primitive

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var withModeInput = []struct {
	input    Mode
	expected []string
}{
	{ModeCombo, []string{"-m", "0"}},
	{ModeTriangle, []string{"-m", "1"}},
}

func TestWithMode(t *testing.T) {
	for _, item := range withModeInput {
		valFunc := WithMode(item.input)
		val := valFunc()
		assert.Equalf(t, item.expected, val, "Expected %v got %v", item.expected, val)
	}
}

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

var primitiveInput = []struct {
	input   []string
	isError bool
}{
	{[]string{"../temp/test_input.png", "../temp/test_output.png"}, false},
	{[]string{"file_not_exist", "file_not_create"}, true},
}

func TestPrimitive(t *testing.T) {
	args := []string{"2"}
	for _, item := range primitiveInput {
		_, err := primitive(item.input[0], item.input[1], 10, args...)
		assert.Equalf(t, item.isError, err != nil, "Expected %v got %v", item.isError, err != nil)
	}

}

var transformInput = []struct {
	testCase string
	input    []string
	isError  bool
}{
	{"tc1", []string{"../temp/test_input.png", "png"}, false},
	{"tc2", []string{"invalid_img_input.png", "png"}, true},
	{"tc3", []string{"invalid_img_input.png", "/"}, true},
	{"tc4", []string{"invalid_img_input.png", "png"}, true},
	{"tc5", []string{"../temp/test_input.png", "png"}, true},
	{"tc6", []string{"invalid_img_input.png", "png"}, true},
}

func TestTransform(t *testing.T) {
	tempIoCopyVar := ioCopyVar
	tempCreateTempFileVar := createTempFileVar
	for _, item := range transformInput {
		imgFile, err := os.OpenFile(item.input[0], os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			t.Error("Failed to Open the file")
		}
		if item.testCase == "tc4" {
			ioCopyVar = func(dst io.Writer, src io.Reader) (written int64, err error) {
				return -1, errors.New("IoCopy failed")
			}
		}
		if item.testCase == "tc5" {
			ioCopyVar = func(dst io.Writer, src io.Reader) (written int64, err error) {
				ioCopyVar = func(dst io.Writer, src io.Reader) (written int64, err error) {
					return -1, errors.New("Io copy failed")
				}
				return io.Copy(dst, src)
			}
		}
		if item.testCase == "tc6" {
			createTempFileVar = func(name, ext string) (*os.File, error) {
				createTempFileVar = func(name, ext string) (*os.File, error) {
					return nil, errors.New("Temp file failed to create")
				}
				return createTempFile(name, ext)
			}
		}
		_, err = Transform(imgFile, item.input[1], 10, WithMode(ModeCombo))
		assert.Equalf(t, item.isError, err != nil, "Expected %v got %v", item.isError, err != nil)
		ioCopyVar = tempIoCopyVar
		createTempFileVar = tempCreateTempFileVar
		imgFile.Close()

		if item.input[0] == "invalid_img_input.png" {
			os.Remove("invalid_img_input.png")
		}
	}
}

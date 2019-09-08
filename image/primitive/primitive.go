package primitive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

//Mode is of type int.
//It is used for primitive modes.
type Mode int

//const below are of type Mode which is basically a number map to a Mode name.
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	Modepolygon
)

var (
	ioCopyVar         = io.Copy
	createTempFileVar = createTempFile
)

//WithMode take a Mode type argument and
//returns a string array containing -m ModeName
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

//Transform takes a io reader, file extension, no. of shapes and other options
//to perform on primitive. It will return io reader having the contents after processing by primitive
//with the options provided.
func Transform(image io.Reader, ext string, shapes int, options ...func() []string) (io.Reader, error) {
	var args []string
	for _, option := range options {
		args = append(args, option()...)
	}
	inFile, err := createTempFileVar("in_", ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(inFile.Name())
	outFile, err := createTempFileVar("out_", ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(outFile.Name())
	_, err = ioCopyVar(inFile, image)
	if err != nil {
		log.Println("Failed to copy:", err)
		return nil, err
	}
	_, err = primitive(inFile.Name(), outFile.Name(), shapes, args...)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(nil)
	_, err = ioCopyVar(b, outFile)
	if err != nil {
		log.Println("Failed to copy:", err)
		return nil, err
	}
	return b, nil
}

//primitive take an input, output file, no. of shapes, different arguments.
//It builds a string having all the options provided in the argument
//and execute the primitive with those optiions.
func primitive(inputFile, outputFile string, shapes int, args ...string) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, shapes)
	args = append(strings.Fields(argStr), args...)
	cmd := exec.Command("primitive", args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("CombinedOutput return error:", err)
		return "", err
	}
	return string(b), nil
}

//createTempFile generates a temporary file.
func createTempFile(name, ext string) (*os.File, error) {
	tFile, err := ioutil.TempFile("", name)
	if err != nil {
		log.Println("Failed to create temp file: ", err)
		return nil, err
	}
	defer os.Remove(tFile.Name())
	file, err := os.Create(fmt.Sprintf("%s.%s", tFile.Name(), ext))
	if err != nil {
		log.Println("Failed to create file: ", err)
		return nil, err
	}
	return file, err
}

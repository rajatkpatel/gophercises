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

type Mode int

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

var ioCopyVar = io.Copy

func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

func Transform(image io.Reader, ext string, shapes int, options ...func() []string) (io.Reader, error) {
	var args []string
	for _, option := range options {
		args = append(args, option()...)
	}
	inFile, err := createTempFile("in_", ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(inFile.Name())
	outFile, err := createTempFile("out_", ext)
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

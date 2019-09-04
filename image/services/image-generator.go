package services

import (
	"io"
	"log"

	"github.com/gophercises/image/primitive"
)

var (
	TempName  = ""
	ioCopyVar = io.Copy
)

//generateImages take io reader, file extension and different options of type generateOptions as argument.
//It calls the generateImage function and return a array of string contating the fileNames returned by generateImage.
func generateImages(rs io.ReadSeeker, ext string, options ...generateOptions) ([]string, error) {
	var ret []string
	for _, option := range options {
		_, err := rs.Seek(0, 0)
		if err != nil {
			log.Println("Failed to seek: ", err)
			return nil, err
		}
		fileName, err := generateImage(rs, ext, option.N, option.M)
		if err != nil {
			log.Println("Failed to generate image:", err)
			return nil, err
		}
		ret = append(ret, fileName)
	}
	return ret, nil

}

//generateImage take io reader, file extension, no. of shapes and primitive Mode as argument.
//It transform the io reader with the options provided and return the new fileName.
func generateImage(r io.Reader, ext string, shapes int, mode primitive.Mode) (string, error) {
	output, err := primitive.Transform(r, ext, shapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}
	outputFile, err := createTempFile(TempName, ext)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	_, err = ioCopyVar(outputFile, output)
	if err != nil {
		log.Println("Failed to copy:", err)
		return "", err
	}
	return outputFile.Name(), nil
}

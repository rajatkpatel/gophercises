package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gophercises/image/primitive"
)

var createTempFileVar = createTempFile

type generateOptions struct {
	N int
	M primitive.Mode
}

//IndexHandler is the handler for the index page.
//It will display the html layout only.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	html := `<html>
				<body>
					<form action="/upload" method="post" enctype="multipart/form-data">
						<input type="file" name="image">
						<button type="submit">Upload Image</button>
					</form>
				</body>
			</html>`
	fmt.Fprint(w, html)
}

//ModifyHandler is the handler for image modify.
//It will display different images after the modification has been done.
func ModifyHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./img/" + filepath.Base(r.URL.Path))
	if err != nil {
		log.Println("Failed to open file: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	ext := filepath.Ext(file.Name())[1:]
	modeStr := r.FormValue("mode")
	if modeStr == "" {
		renderModeChoices(w, r, file, ext)
		return
	}
	mode, err := strconv.Atoi(modeStr)
	if err != nil {
		log.Println("String to integer conversion failed: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nStr := r.FormValue("n")
	if nStr == "" {
		renderShapeChoices(w, r, file, ext, primitive.Mode(mode))
		return
	}
	shape, err := strconv.Atoi(nStr)
	if err != nil {
		log.Println("String to integer conversion failed: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = shape
	http.Redirect(w, r, "/img/"+filepath.Base(file.Name()), http.StatusFound)
}

//UploadHandler is the handler for uploading a file.
//It will redirect to modify after successful upload of a file.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("No file found for the key provided: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)[1:]
	onDiskFile, err := createTempFileVar("", ext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer onDiskFile.Close()

	_, err = ioCopyVar(onDiskFile, file)
	if err != nil {
		log.Println("Failed to copy:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/modify/"+filepath.Base(onDiskFile.Name()), http.StatusFound)
}

//createTempFile generates a temporary file.
func createTempFile(name, ext string) (*os.File, error) {
	tFile, err := ioutil.TempFile("./img/", name)
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

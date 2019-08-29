package services

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gophercises/image/primitive"
)

func renderShapeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string, mode primitive.Mode) {
	options := []generateOptions{
		{N: 10, M: mode},
		{N: 20, M: mode},
		{N: 30, M: mode},
		{N: 40, M: mode},
	}

	imgs, err := generateImages(rs, ext, options...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
					{{range .}}
						<a href = "/modify/{{.Name}}?mode={{.Mode}}&n={{.Shape}}">
						<img style="width: 25%;" src = "/img/{{.Name}}">
						</a>
					{{end}}
				</body></html>`

	tpl := template.Must(template.New("").Parse(html))

	type dataStruct struct {
		Name  string
		Mode  primitive.Mode
		Shape int
	}

	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name:  filepath.Base(img),
			Mode:  options[i].M,
			Shape: options[i].N,
		})
	}

	err = tpl.Execute(w, data)
	if err != nil {
		log.Println("Template execution failed: ", err)
	}
}

func renderModeChoices(w http.ResponseWriter, r *http.Request, rs io.ReadSeeker, ext string) {
	options := []generateOptions{
		{N: 20, M: primitive.ModeCircle},
		{N: 20, M: primitive.ModeCombo},
		{N: 20, M: primitive.ModeRect},
		{N: 20, M: primitive.ModeEllipse},
	}

	imgs, err := generateImages(rs, ext, options...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := `<html><body>
					{{range .}}
						<a href = "/modify/{{.Name}}?mode={{.Mode}}">
						<img style="width: 25%;" src = "/img/{{.Name}}">
						</a>
					{{end}}
				</body></html>`

	tpl := template.Must(template.New("").Parse(html))

	type dataStruct struct {
		Name string
		Mode primitive.Mode
	}

	var data []dataStruct
	for i, img := range imgs {
		data = append(data, dataStruct{
			Name: filepath.Base(img),
			Mode: options[i].M,
		})
	}

	err = tpl.Execute(w, data)
	if err != nil {
		log.Println("Template execution failed: ", err)
	}
}

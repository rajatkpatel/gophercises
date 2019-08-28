package services

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	IoCopy        = io.Copy
	lexer         = lexers.Get("go")
	LexerTokenise = lexer.Tokenise
)

//PanicHandler take two parameters i.e. http ResponseWriter and http Request.
//It cause the code to panic.
func PanicHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Panic Called")
	panic("Oh no! Code Panics")
}

//SourceCodeHandler take two parameters i.e. http ResponseWriter and http Request.
//It will show the source code file with lines highlighting
//where file name and line number are fetched from the URL.
func SourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	lineNo := r.FormValue("line")
	line, err := strconv.Atoi(lineNo)
	if err != nil {
		line = -1
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := bytes.NewBuffer(nil)
	_, err = IoCopy(b, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var lines [][2]int
	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}

	iterator, err := LexerTokenise(nil, b.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	style := styles.Get("github")
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(), html.LineNumbersInTable(), html.HighlightLines(lines))
	w.Header().Set("Content-Type", "text/html")
	formatter.Format(w, style, iterator)

	//err = quick.Highlight(w, b.String(), "go", "html", "monokai")
}

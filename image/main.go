package main

import (
	"log"
	"net/http"

	"github.com/gophercises/image/services"
)

var (
	listenAndServe = http.ListenAndServe
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", services.IndexHandler)
	mux.HandleFunc("/modify/", services.ModifyHandler)
	mux.HandleFunc("/upload", services.UploadHandler)

	fileServer := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img", fileServer))

	log.Fatal(listenAndServe(":5000", mux))
}

func init() {
	log.SetPrefix("LOG: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

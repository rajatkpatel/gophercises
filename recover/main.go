package main

import (
	"log"
	"net/http"

	"github.com/gophercises/recover/middleware"
	"github.com/gophercises/recover/services"
)

var (
	listenAndServe = http.ListenAndServe
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug", services.SourceCodeHandler)
	mux.HandleFunc("/panic", services.PanicHandler)
	log.Fatal(listenAndServe(":5000", middleware.RecoverMw(mux)))
}

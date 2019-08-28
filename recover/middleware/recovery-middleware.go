package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gophercises/recover/hyperlink"
)

//RecoverMw takes a parameter of http handler type and it will return a http handler function as ouput.
//It will help to recover from panic and print a stacktrace where the panic occured.
func RecoverMw(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, hyperlink.CreateLinks(string(stack)))
			}
		}()
		handler.ServeHTTP(w, r)
	}
}

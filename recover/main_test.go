package main

import (
	"testing"
	"time"
)

func TestM(t *testing.T) {
	// tempListenAndServe := listenAndServe
	// defer func() {
	// 	listenAndServe = tempListenAndServe
	// }()
	// listenAndServe = func(addr string, handler http.Handler) error {
	// 	return errors.New("Failed")
	// }
	go main()
	time.Sleep(1 * time.Second)
}

package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestM(t *testing.T) {
	// go main()
	// time.Sleep(1 * time.Second)
	tempListenAndServe := listenAndServe
	defer func() {
		listenAndServe = tempListenAndServe
	}()
	listenAndServe = func(addr string, handler http.Handler) error {
		panic("Failed")
	}
	assert.PanicsWithValuef(t, "Failed", main, "Expected %v got %v", "Failed", main)
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gophercises/recover/services"
)

func TestRecoveryMw(t *testing.T) {
	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Error("Error occured while testing: ", err)
	}
	rr := httptest.NewRecorder()
	panicHandler := http.HandlerFunc(services.PanicHandler)
	handler := http.HandlerFunc(RecoverMw(panicHandler))
	handler.ServeHTTP(rr, req)
}

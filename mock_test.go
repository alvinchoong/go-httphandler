package httphandler_test

import (
	"net/http"
)

// mockResponder is a mock implementation of the Responder interface.
type mockResponder struct {
	StatusCode int
	Body       string
}

// Respond writes the preset status code and body to the ResponseWriter.
func (mr *mockResponder) Respond(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(mr.StatusCode)
	w.Write([]byte(mr.Body))
}

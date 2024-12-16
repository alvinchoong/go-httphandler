package httphandler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/jsonresp"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var testUser = &User{
	ID:   1,
	Name: "John Doe",
}

// Standard Go HTTP handler
func standardHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testUser)
}

// go-httphandler handler
func customHandler() http.HandlerFunc {
	return httphandler.Handle(func(r *http.Request) httphandler.Responder {
		return jsonresp.Success(testUser)
	})
}

func BenchmarkJSONResponse(b *testing.B) {
	b.Run("Go/StandardHTTPHandler", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest(http.MethodGet, "/user", nil)
			w := httptest.NewRecorder()
			standardHandler(w, r)
		}
	})

	b.Run("HTTPHandler/JSONResponse", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest(http.MethodGet, "/user", nil)
			w := httptest.NewRecorder()
			customHandler()(w, r)
		}
	})
}

// Standard Go HTTP handler + json decode
func standardHandlerAndJSONDecode(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
}

// go-httphandler HandleWithInput
func customHandlerWithInput() http.HandlerFunc {
	return httphandler.HandleWithInput(func(r *http.Request, input User) httphandler.Responder {
		return nil
	})
}

func BenchmarkJSONRequest(b *testing.B) {
	inputJSON, _ := json.Marshal(testUser)

	b.Run("Go/StandardHTTPHandler/JSONDecode", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(inputJSON))
			w := httptest.NewRecorder()
			standardHandlerAndJSONDecode(w, r)
		}
	})

	b.Run("HTTPHandler/JSONBodyDecode", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(inputJSON))
			w := httptest.NewRecorder()
			customHandlerWithInput()(w, r)
		}
	})
}

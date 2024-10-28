package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Responder defines how to respond to HTTP requests.
type Responder interface {
	Respond(w http.ResponseWriter, r *http.Request)
}

// RequestHandler handles an HTTP request and returns a Responder.
type RequestHandler func(r *http.Request) Responder

// Handle converts a RequestHandler to an http.HandlerFunc.
func Handle(handler RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := handler(r)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// RequestDecodeFunc defines how to decode an HTTP request.
type RequestDecodeFunc[T any] func(r *http.Request) (T, error)

// RequestHandlerWithInput handles an HTTP request with decoded input and returns a Responder.
type RequestHandlerWithInput[T any] func(r *http.Request, input T) Responder

// HandleWithInput converts a RequestHandlerWithInput to an http.HandlerFunc.
type handleWithInput[T any] struct {
	decodeFunc RequestDecodeFunc[T]
	handler    RequestHandlerWithInput[T]
}

// HandleWithInput converts a RequestHandlerWithInput to an http.HandlerFunc.
func HandleWithInput[T any](handler RequestHandlerWithInput[T], opts ...func(*handleWithInput[T])) http.HandlerFunc {
	h := &handleWithInput[T]{
		decodeFunc: JSONBodyDecode[T],
		handler:    handler,
	}
	for _, opt := range opts {
		opt(h)
	}

	return h.ServeHTTP
}

// ServeHTTP implements the http.Handler interface.
func (h *handleWithInput[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	input, err := h.decodeFunc(r)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	res := h.handler(r, input)
	if res == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	res.Respond(w, r)
}

// WithDecodeFunc sets the decode function for the handler.
func WithDecodeFunc[T any](decodeFunc RequestDecodeFunc[T]) func(*handleWithInput[T]) {
	return func(h *handleWithInput[T]) {
		h.decodeFunc = decodeFunc
	}
}

var ErrJSONDecode = errors.New("fail to decode json")

func JSONBodyDecode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("%w: %w", ErrJSONDecode, err)
	}

	return v, nil
}

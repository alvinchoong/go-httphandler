package jsonresp

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/alvinchoong/go-httphandler"
)

// Ensure successResponder implements Responder.
var _ httphandler.Responder = (*successResponder[any])(nil)

// Success creates a new successResponder with the provided data and a default status code of 200 OK.
func Success[T any](data *T) *successResponder[T] {
	return &successResponder[T]{
		statusCode: http.StatusOK,
		data:       data,
	}
}

// successResponder handles successful JSON HTTP responses.
type successResponder[T any] struct {
	logger     httphandler.Logger
	header     http.Header
	statusCode int
	cookies    []*http.Cookie
	data       *T
}

// Respond sends the JSON response with custom headers, cookies and status code.
func (res *successResponder[T]) Respond(w http.ResponseWriter, _ *http.Request) {
	// Set cookies.
	for _, cookie := range res.cookies {
		http.SetCookie(w, cookie)
	}

	// Add custom headers.
	for key, values := range res.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write the JSON response.
	b := writeJSON(w, res.data, res.statusCode, res.logger)
	httphandler.LogResponse(res.logger, res.statusCode, "response_body", b)
}

// WithLogger sets the logger for the responder.
func (res *successResponder[T]) WithLogger(logger httphandler.Logger) *successResponder[T] {
	res.logger = logger
	return res
}

// WithStatus sets a custom HTTP status code for the response.
func (res *successResponder[T]) WithStatus(status int) *successResponder[T] {
	res.statusCode = status
	return res
}

// WithHeader adds a custom header to the response.
func (res *successResponder[T]) WithHeader(key, value string) *successResponder[T] {
	if res.header == nil {
		res.header = http.Header{}
	}
	res.header.Add(key, value)
	return res
}

// WithCookie adds a cookie to the response.
func (res *successResponder[T]) WithCookie(cookie *http.Cookie) *successResponder[T] {
	res.cookies = append(res.cookies, cookie)
	return res
}

// writeJSON encodes the data as JSON and writes it to the ResponseWriter with the specified status code.
// If encoding fails, it responds with a 500 Internal Server Error.
func writeJSON(w http.ResponseWriter, v any, status int, logger httphandler.Logger) []byte {
	w.Header().Set("Content-Type", "application/json")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		httphandler.WriteInternalServerError(w, logger, err, "data", v)
		return nil
	}

	w.WriteHeader(status)
	if _, err := w.Write(buf.Bytes()); err != nil {
		httphandler.WriteInternalServerError(w, logger, err, "response_body", buf.String())
		return nil
	}

	return buf.Bytes()
}

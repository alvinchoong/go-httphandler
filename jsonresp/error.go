package jsonresp

import (
	"net/http"

	"github.com/alvinchoong/go-httphandler"
)

// Ensure errorResponder implements Responder.
var _ httphandler.Responder = (*errorResponder[any])(nil)

// Error creates a standardized error response with the specified error message and HTTP status code.
// The 'err' parameter can be used for internal logging.
func Error[T any](err error, errData T, code int) *errorResponder[T] {
	return &errorResponder[T]{
		statusCode: code,
		errData:    errData,
		err:        err,
	}
}

// InternalServerError creates a standardized internal server error response.
// The 'err' parameter can be used for internal logging.
func InternalServerError(err error) *errorResponder[string] {
	return &errorResponder[string]{
		statusCode: http.StatusInternalServerError,
		errData:    "Internal Server Error",
		err:        err,
	}
}

// errorResponder handles error JSON HTTP responses.
type errorResponder[T any] struct {
	logger     httphandler.Logger
	header     http.Header
	statusCode int
	cookies    []*http.Cookie
	errData    T
	err        error
}

// Respond sends the JSON error response with custom headers, cookies, and status code.
func (res *errorResponder[T]) Respond(w http.ResponseWriter, _ *http.Request) {
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

	// Write the error JSON response.
	writeJSON(w, map[string]T{"error": res.errData}, res.statusCode, res.logger)
	httphandler.LogRequestError(res.logger, res.err)
}

// WithLogger sets the logger for the responder.
func (res *errorResponder[T]) WithLogger(logger httphandler.Logger) *errorResponder[T] {
	res.logger = logger
	return res
}

// WithHeader adds a custom header to the response.
func (res *errorResponder[T]) WithHeader(key, value string) *errorResponder[T] {
	if res.header == nil {
		res.header = http.Header{}
	}
	res.header.Add(key, value)
	return res
}

// WithCookie adds a cookie to the response.
func (res *errorResponder[T]) WithCookie(cookie *http.Cookie) *errorResponder[T] {
	res.cookies = append(res.cookies, cookie)
	return res
}

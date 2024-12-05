package jsonresp

import (
	"net/http"

	"github.com/alvinchoong/go-httphandler"
)

// Ensure errorResponder implements Responder.
var _ httphandler.Responder = (*errorResponder)(nil)

// Error creates a standardized error response with the specified error message and HTTP status code.
// The 'err' parameter can be used for internal logging.
func Error(err error, message string, code int) *errorResponder {
	return &errorResponder{
		statusCode: code,
		errMessage: message,
		err:        err,
	}
}

// InternalServerError creates a standardized internal server error response.
// The 'err' parameter can be used for internal logging.
func InternalServerError(err error) *errorResponder {
	return &errorResponder{
		statusCode: http.StatusInternalServerError,
		errMessage: "Internal Server Error",
		err:        err,
	}
}

// errorResponder handles error JSON HTTP responses.
type errorResponder struct {
	logger     httphandler.Logger
	header     http.Header
	statusCode int
	cookies    []*http.Cookie
	errMessage string
	err        error
}

// Respond sends the JSON error response with custom headers, cookies, and status code.
func (res *errorResponder) Respond(w http.ResponseWriter, _ *http.Request) {
	// Add custom headers.
	for key, values := range res.header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set cookies.
	for _, cookie := range res.cookies {
		http.SetCookie(w, cookie)
	}

	// Write the error JSON response.
	if b := writeJSON(w, map[string]string{"error": res.errMessage}, res.statusCode, res.logger); b != nil {
		if res.logger != nil {
			res.logger.Error("Sent error HTTP response", "error", res.err)
			res.logger.Info("Sent HTTP response",
				"status_code", res.statusCode,
				"response_body", string(b),
			)
		}
	}
}

// WithLogger sets the logger for the responder.
func (res *errorResponder) WithLogger(logger httphandler.Logger) *errorResponder {
	res.logger = logger
	return res
}

// WithHeader adds a custom header to the response.
func (res *errorResponder) WithHeader(key, value string) *errorResponder {
	if res.header == nil {
		res.header = http.Header{}
	}
	res.header.Add(key, value)
	return res
}

// WithCookie adds a cookie to the response.
func (res *errorResponder) WithCookie(cookie *http.Cookie) *errorResponder {
	res.cookies = append(res.cookies, cookie)
	return res
}

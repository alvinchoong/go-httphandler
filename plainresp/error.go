package plainresp

import (
	"net/http"

	"github.com/alvinchoong/go-httphandler"
)

// Ensure errorResponder implements Responder.
var _ httphandler.Responder = (*errorResponder)(nil)

// errorResponder represents an error response with a message.
type errorResponder struct {
	logger     httphandler.Logger
	header     http.Header
	statusCode int
	cookies    []*http.Cookie
	body       string
	err        error
}

// Error creates a new errorResponder with the provided error, message, and status code.
// The 'err' parameter can be used for internal logging.
func Error(err error, body string, code int) *errorResponder {
	return &errorResponder{
		statusCode: code,
		body:       body,
		err:        err,
	}
}

// InternalServerError creates a standardized internal server error response.
// The 'err' parameter can be used for internal logging.
func InternalServerError(err error) *errorResponder {
	return &errorResponder{
		statusCode: http.StatusInternalServerError,
		body:       "Internal Server Error",
		err:        err,
	}
}

// Respond sends the response with custom headers, cookies and status code.
func (res *errorResponder) Respond(w http.ResponseWriter, _ *http.Request) {
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

	// Set response body and status code.
	http.Error(w, res.body, res.statusCode)
	httphandler.LogRequestError(res.logger, res.err)
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

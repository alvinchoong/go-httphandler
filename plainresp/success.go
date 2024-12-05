package plainresp

import (
	"net/http"

	"github.com/alvinchoong/go-httphandler"
)

// Ensure successResponder implements Responder.
var _ httphandler.Responder = (*successResponder)(nil)

// successResponder manages successful HTTP responses.
type successResponder struct {
	logger     httphandler.Logger
	header     http.Header
	statusCode int
	cookies    []*http.Cookie
	body       string
}

// Success creates a new successResponder with data and a 200 OK status.
func Success(data string) *successResponder {
	return &successResponder{
		statusCode: http.StatusOK,
		body:       data,
	}
}

// Respond sends the response with custom headers, cookies and status code.
func (res *successResponder) Respond(w http.ResponseWriter, r *http.Request) {
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

	// Set response body and status code.
	w.WriteHeader(res.statusCode)
	if _, err := w.Write([]byte(res.body)); err != nil {
		if res.logger != nil {
			res.logger.Error("Failed to write HTTP response",
				"error", err,
			)
		}
		return
	}

	if res.logger != nil {
		res.logger.Info("Sent HTTP response",
			"status_code", res.statusCode,
			"response_body", res.body,
		)
	}
}

// WithLogger sets the logger for the responder.
func (res *successResponder) WithLogger(logger httphandler.Logger) *successResponder {
	res.logger = logger
	return res
}

// WithHeader adds a custom header to the response.
func (res *successResponder) WithHeader(key, value string) *successResponder {
	if res.header == nil {
		res.header = http.Header{}
	}
	res.header.Add(key, value)
	return res
}

// WithStatus sets a custom HTTP status code for the response.
func (res *successResponder) WithStatus(status int) *successResponder {
	res.statusCode = status
	return res
}

// WithCookie adds a cookie to the response.
func (res *successResponder) WithCookie(cookie *http.Cookie) *successResponder {
	res.cookies = append(res.cookies, cookie)
	return res
}

package templresp

import (
	"context"
	"io"
	"net/http"

	"github.com/alvinchoong/go-httphandler"
)

// Ensure responder implements Responder.
var _ httphandler.Responder = (*contentResponder)(nil)

// contentResponder manages successful HTTP responses.
type contentResponder struct {
	logger     httphandler.Logger
	header     http.Header
	statusCode int
	cookies    []*http.Cookie
	component  Component
}

// Component is the interface to integrate with Templ components
type Component interface {
	Render(ctx context.Context, w io.Writer) error
}

// Content creates a new contentResponder with the provided component and a default status code of 200 OK.
func Content(component Component) *contentResponder {
	return &contentResponder{
		statusCode: http.StatusOK,
		component:  component,
	}
}

// Respond sends the content response with custom headers, cookies and status code.
func (res *contentResponder) Respond(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set response body and status code.
	w.WriteHeader(res.statusCode)
	if err := res.component.Render(r.Context(), w); err != nil {
		httphandler.WriteInternalServerError(w, res.logger, err)
		return
	}

	httphandler.LogResponse(res.logger, res.statusCode)
}

// WithLogger sets the logger for the responder.
func (res *contentResponder) WithLogger(logger httphandler.Logger) *contentResponder {
	res.logger = logger
	return res
}

// WithStatus sets a custom HTTP status code for the response.
func (res *contentResponder) WithStatus(status int) *contentResponder {
	res.statusCode = status
	return res
}

// WithHeader adds a custom header to the response.
func (res *contentResponder) WithHeader(key, value string) *contentResponder {
	if res.header == nil {
		res.header = http.Header{}
	}
	res.header.Add(key, value)
	return res
}

// WithCookie adds a cookie to the response.
func (res *contentResponder) WithCookie(cookie *http.Cookie) *contentResponder {
	res.cookies = append(res.cookies, cookie)
	return res
}

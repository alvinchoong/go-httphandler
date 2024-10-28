package httphandler

import "net/http"

// Ensure redirectResponder implements Responder.
var _ Responder = (*redirectResponder)(nil)

// Redirect creates a new redirectResponder with the specified URL and status code.
func Redirect(url string, code int) *redirectResponder {
	return &redirectResponder{
		statusCode: code,
		url:        url,
	}
}

// redirectResponder handles HTTP redirects.
type redirectResponder struct {
	logger     Logger
	header     http.Header
	statusCode int
	cookies    []*http.Cookie
	url        string
}

// Respond sents an HTTP redirect with custom headers, cookies, and status code.
func (res *redirectResponder) Respond(w http.ResponseWriter, r *http.Request) {
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

	// Redirect to the specified URL.
	http.Redirect(w, r, res.url, res.statusCode)
	if res.logger != nil {
		res.logger.Info("Sent HTTP redirect",
			"status_code", res.statusCode,
			"redirect_url", res.url,
		)
	}
}

// WithLogger sets the logger for the responder.
// This allows the redirectResponder to log any relevant information.
func (res *redirectResponder) WithLogger(logger Logger) *redirectResponder {
	res.logger = logger
	return res
}

// WithHeader adds a header to the response.
func (res *redirectResponder) WithHeader(key, value string) *redirectResponder {
	if res.header == nil {
		res.header = http.Header{}
	}
	res.header.Add(key, value)
	return res
}

// WithCookie adds a cookie to the response.
func (res *redirectResponder) WithCookie(cookie *http.Cookie) *redirectResponder {
	res.cookies = append(res.cookies, cookie)
	return res
}

package downloadresp

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/alvinchoong/go-httphandler"
)

// Ensure responder implements Responder.
var _ httphandler.Responder = (*fileResponder)(nil)

// fileResponder handles file response that can be returned from an HTTP handler.
type fileResponder struct {
	logger      httphandler.Logger
	header      http.Header
	cookies     []*http.Cookie
	reader      io.Reader
	filename    string
	disposition string
}

// Attachment returns a responder that can be used to send a file as an attachment.
func Attachment(reader io.Reader, filename string) *fileResponder {
	return &fileResponder{
		reader:      reader,
		filename:    filename,
		disposition: "attachment",
	}
}

// Inline returns a responder that can be used to send a file as inline.
func Inline(reader io.Reader, filename string) *fileResponder {
	return &fileResponder{
		reader:      reader,
		filename:    filename,
		disposition: "inline",
	}
}

// Respond sends the response with custom headers.
func (res *fileResponder) Respond(w http.ResponseWriter, _ *http.Request) {
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

	// If the Content-Type header is not set, set it to the appropriate MIME type based on the file extension.
	if res.header.Get("Content-Type") == "" {
		contentType := "application/octet-stream"
		if s := mime.TypeByExtension(filepath.Ext(res.filename)); s != "" {
			contentType = s
		}
		w.Header().Set("Content-Type", contentType)
	}

	// Set the Content-Disposition header
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`%s; filename="%s"`, res.disposition, res.filename))

	if _, err := io.Copy(w, res.reader); err != nil {
		if res.logger != nil {
			res.logger.Error("Failed to write HTTP response",
				"error", err,
			)
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if res.logger != nil {
		res.logger.Info("Sent HTTP response",
			"status_code", http.StatusOK,
			"filename", res.filename,
		)
	}
}

// WithHeader adds a custom header to the response.
func (res *fileResponder) WithHeader(key, value string) *fileResponder {
	if res.header == nil {
		res.header = http.Header{}
	}
	res.header.Add(key, value)
	return res
}

// WithContentType sets the Content-Type header.
func (res *fileResponder) WithContentType(contentType string) *fileResponder {
	return res.WithHeader("Content-Type", contentType)
}

// WithLogger sets the logger for the responder.
func (res *fileResponder) WithLogger(logger httphandler.Logger) *fileResponder {
	res.logger = logger
	return res
}

// WithCookie adds a cookie to the response.
func (res *fileResponder) WithCookie(cookie *http.Cookie) *fileResponder {
	res.cookies = append(res.cookies, cookie)
	return res
}

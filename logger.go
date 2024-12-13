package httphandler

import "net/http"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// WriteInternalServerError writes an HTTP 500 Internal Server Error response.
// If a logger is provided, it will also log the error.
func WriteInternalServerError(w http.ResponseWriter, logger Logger, err error, args ...any) {
	if logger != nil {
		logger.Error("Failed to write HTTP response",
			append([]any{"error", err}, args...)...,
		)
	}
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// LogResponse logs the response status if a logger is provided.
func LogResponse(logger Logger, status int, args ...any) {
	if logger == nil {
		return
	}

	logger.Info("Sent HTTP response",
		append([]any{"status_code", status}, args...)...,
	)
}

// LogRequestError logs an error if a logger is provided.
func LogRequestError(logger Logger, err error, args ...any) {
	if logger == nil {
		return
	}

	logger.Error("Error handling request",
		append([]any{"error", err}, args...)...,
	)
}

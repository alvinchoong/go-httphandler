package httphandler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ========== Header Decoders ==========

// HeaderValue returns a decoder that extracts a specific header value
func HeaderValue(name string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		value := r.Header.Get(name)
		if value == "" {
			return "", fmt.Errorf("header %q not found", name)
		}
		return value, nil
	}
}

// OptionalHeaderValue returns a decoder that extracts a header value if present
func OptionalHeaderValue(name string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		return r.Header.Get(name), nil
	}
}

// BearerToken returns a decoder that extracts a bearer token from the Authorization header
func BearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("invalid Authorization header format, expected 'Bearer TOKEN'")
	}

	return strings.TrimPrefix(auth, "Bearer "), nil
}

// OptionalBearerToken returns a decoder that extracts a bearer token if present
func OptionalBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return "", nil
	}

	return strings.TrimPrefix(auth, "Bearer "), nil
}

// BasicAuth returns a decoder that extracts username and password from Basic Auth
func BasicAuth(r *http.Request) (struct {
	Username string
	Password string
}, error,
) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return struct {
			Username string
			Password string
		}{}, fmt.Errorf("missing or invalid Basic Auth")
	}

	return struct {
		Username string
		Password string
	}{
		Username: username,
		Password: password,
	}, nil
}

// ========== Query Parameter Decoders ==========

// QueryParam returns a decoder that extracts a query parameter value
func QueryParam(name string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		value := r.URL.Query().Get(name)
		if value == "" {
			return "", fmt.Errorf("query parameter %q not found", name)
		}
		return value, nil
	}
}

// OptionalQueryParam returns a decoder that extracts a query parameter value if present
func OptionalQueryParam(name string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		return r.URL.Query().Get(name), nil
	}
}

// QueryParams returns a decoder that extracts all query parameters
func QueryParams(r *http.Request) (url.Values, error) {
	return r.URL.Query(), nil
}

// ParsedQueryParam returns a decoder that extracts and parses a query parameter to a specific type
func ParsedQueryParam[T any](name string, parser func(string) (T, error)) func(r *http.Request) (T, error) {
	return func(r *http.Request) (T, error) {
		var zero T
		value := r.URL.Query().Get(name)
		if value == "" {
			return zero, fmt.Errorf("query parameter %q not found", name)
		}

		parsed, err := parser(value)
		if err != nil {
			return zero, fmt.Errorf("failed to parse query parameter %q: %w", name, err)
		}

		return parsed, nil
	}
}

// IntQueryParam returns a decoder that extracts a query parameter as an integer
func IntQueryParam(name string) func(r *http.Request) (int, error) {
	return ParsedQueryParam(name, strconv.Atoi)
}

// FloatQueryParam returns a decoder that extracts a query parameter as a float64
func FloatQueryParam(name string) func(r *http.Request) (float64, error) {
	return ParsedQueryParam(name, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

// BoolQueryParam returns a decoder that extracts a query parameter as a boolean
func BoolQueryParam(name string) func(r *http.Request) (bool, error) {
	return ParsedQueryParam(name, strconv.ParseBool)
}

// ========== Path Parameter Decoders ==========

// PathParam returns a decoder that extracts a path parameter value
// This works with Go 1.22's r.PathValue() method
func PathParam(name string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		value := r.PathValue(name)
		if value == "" {
			return "", fmt.Errorf("path parameter %q not found", name)
		}
		return value, nil
	}
}

// ParsedPathParam returns a decoder that extracts and parses a path parameter to a specific type
func ParsedPathParam[T any](name string, parser func(string) (T, error)) func(r *http.Request) (T, error) {
	return func(r *http.Request) (T, error) {
		var zero T
		value := r.PathValue(name)
		if value == "" {
			return zero, fmt.Errorf("path parameter %q not found", name)
		}

		parsed, err := parser(value)
		if err != nil {
			return zero, fmt.Errorf("failed to parse path parameter %q: %w", name, err)
		}

		return parsed, nil
	}
}

// IntPathParam returns a decoder that extracts a path parameter as an integer
func IntPathParam(name string) func(r *http.Request) (int, error) {
	return ParsedPathParam(name, strconv.Atoi)
}

// ========== Form Data Decoders ==========

// FormValue returns a decoder that extracts a form value
func FormValue(name string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		// Parse form data if not already parsed
		if err := r.ParseForm(); err != nil {
			return "", fmt.Errorf("failed to parse form: %w", err)
		}

		value := r.FormValue(name)
		if value == "" {
			return "", fmt.Errorf("form value %q not found", name)
		}
		return value, nil
	}
}

// OptionalFormValue returns a decoder that extracts a form value if present
func OptionalFormValue(name string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		// Parse form data if not already parsed
		if err := r.ParseForm(); err != nil {
			return "", fmt.Errorf("failed to parse form: %w", err)
		}

		return r.FormValue(name), nil
	}
}

// FormFile returns a decoder that extracts a file from a multipart form
type MultipartFile struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func FormFile(name string) func(r *http.Request) (*MultipartFile, error) {
	return func(r *http.Request) (*MultipartFile, error) {
		// Parse multipart form with default max memory
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return nil, fmt.Errorf("failed to parse multipart form: %w", err)
		}

		file, header, err := r.FormFile(name)
		if err != nil {
			return nil, fmt.Errorf("form file %q not found: %w", name, err)
		}

		return &MultipartFile{
			File:   file,
			Header: header,
		}, nil
	}
}

// ========== Body Decoders ==========

// JSONBody decodes the request body as JSON into a value of type T
func JSONBody[T any]() func(r *http.Request) (T, error) {
	return func(r *http.Request) (T, error) {
		var value T
		if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
			return value, fmt.Errorf("failed to decode JSON: %w", err)
		}
		return value, nil
	}
}

// ========== Composite Decoders ==========

// Combine2 combines two decoders into a single decoder that returns a struct with both results
func Combine2[T1, T2 any](
	decoder1 func(r *http.Request) (T1, error),
	decoder2 func(r *http.Request) (T2, error),
) func(r *http.Request) (struct {
	V1 T1
	V2 T2
}, error) {
	return func(r *http.Request) (struct {
		V1 T1
		V2 T2
	}, error,
	) {
		v1, err := decoder1(r)
		if err != nil {
			return struct {
				V1 T1
				V2 T2
			}{}, fmt.Errorf("decoder1 error: %w", err)
		}

		v2, err := decoder2(r)
		if err != nil {
			return struct {
				V1 T1
				V2 T2
			}{}, fmt.Errorf("decoder2 error: %w", err)
		}

		return struct {
			V1 T1
			V2 T2
		}{
			V1: v1,
			V2: v2,
		}, nil
	}
}

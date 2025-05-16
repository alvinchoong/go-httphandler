package jsonresp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/jsonresp"
)

func TestError_Respond(t *testing.T) {
	t.Parallel()

	type ValidationError struct {
		Field   string   `json:"field"`
		Message string   `json:"message"`
		Details []string `json:"details,omitempty"`
	}

	type ErrorResponse struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Errors  []ValidationError `json:"errors,omitempty"`
	}

	cookie := &http.Cookie{
		Name:  "test-cookie-1",
		Value: "cookie-value-1",
	}

	testCases := []struct {
		desc        string
		given       httphandler.Responder
		wantCode    int
		wantHeaders map[string]string
		wantCookies []*http.Cookie
		wantBody    string
	}{
		{
			desc:        "basic | nil",
			given:       jsonresp.Error[any](errors.New("unexpected error"), nil, http.StatusInternalServerError),
			wantCode:    http.StatusInternalServerError,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `null`,
		},
		{
			desc:        "basic | string",
			given:       jsonresp.Error(errors.New("invalid id"), ptr("Invalid ID provided"), http.StatusBadRequest),
			wantCode:    http.StatusBadRequest,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `"Invalid ID provided"`,
		},
		{
			desc:        "basic | slice",
			given:       jsonresp.Error(errors.New("invalid id"), &[]string{"Required field", "Invalid format"}, http.StatusBadRequest),
			wantCode:    http.StatusBadRequest,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `["Required field","Invalid format"]`,
		},
		{
			desc:        "basic | map",
			given:       jsonresp.Error(errors.New("validation failed"), &map[string]string{"email": "Invalid email format", "password": "Too short"}, http.StatusBadRequest),
			wantCode:    http.StatusBadRequest,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `{"email":"Invalid email format","password":"Too short"}`,
		},
		{
			desc: "basic | struct",
			given: jsonresp.Error(
				errors.New("validation failed"),
				&ErrorResponse{
					Code:    "VALIDATION_ERROR",
					Message: "The request contains invalid parameters",
					Errors: []ValidationError{
						{Field: "email", Message: "Invalid email format", Details: []string{"Must contain @", "Domain must be valid"}},
						{Field: "password", Message: "Password too weak", Details: []string{"Must be at least 8 characters", "Must contain a number"}},
					},
				},
				http.StatusBadRequest,
			),
			wantCode:    http.StatusBadRequest,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `{"code":"VALIDATION_ERROR","message":"The request contains invalid parameters","errors":[{"field":"email","message":"Invalid email format","details":["Must contain @","Domain must be valid"]},{"field":"password","message":"Password too weak","details":["Must be at least 8 characters","Must contain a number"]}]}`,
		},
		{
			desc: "with everything",
			given: jsonresp.Error(
				errors.New("resource not found"),
				&map[string]string{"code": "NOT_FOUND", "message": "The requested resource was not found"},
				http.StatusNotFound,
			).
				WithHeader("X-Request-ID", "abc123").
				WithCookie(cookie),
			wantCode: http.StatusNotFound,
			wantHeaders: map[string]string{
				"X-Request-ID": "abc123",
			},
			wantCookies: []*http.Cookie{cookie},
			wantBody:    `{"code":"NOT_FOUND","message":"The requested resource was not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test-error", nil)

			// When:
			tc.given.Respond(w, r)

			// Then:
			gotCode := w.Code
			if gotCode != tc.wantCode {
				t.Errorf("status code: want %d, got %d", tc.wantCode, gotCode)
			}

			gotHeaders := w.Header()

			// Always check for Content-Type header
			contentType := gotHeaders.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type header: want 'application/json', got '%s'", contentType)
			}

			// Check other expected headers
			for key, wantValue := range tc.wantHeaders {
				gotValue := gotHeaders.Get(key)
				if gotValue != wantValue {
					t.Errorf("header %s: want %s, got %s", key, wantValue, gotValue)
				}
			}

			gotCookies := w.Result().Cookies()
			if len(gotCookies) != len(tc.wantCookies) {
				t.Errorf("cookie count: want %d, got %d", len(tc.wantCookies), len(gotCookies))
			}

			for i, want := range tc.wantCookies {
				got := gotCookies[i]
				if got.Name != want.Name || got.Value != want.Value {
					t.Errorf("cookie %d: want {%s=%s}, got {%s=%s}",
						i, want.Name, want.Value, got.Name, got.Value)
				}
			}

			gotBody := strings.TrimSpace(w.Body.String())
			if gotBody != tc.wantBody {
				t.Errorf("body: want '%s', got '%s'", tc.wantBody, gotBody)
			}
		})
	}
}

func TestInternalServerError_Respond(t *testing.T) {
	t.Parallel()

	cookie := &http.Cookie{
		Name:     "test-cookie-2",
		Value:    "cookie-value-2",
		Path:     "/",
		Domain:   "example.com",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}

	testCases := []struct {
		desc        string
		given       httphandler.Responder
		wantCode    int
		wantHeaders map[string]string
		wantCookies []*http.Cookie
		wantBody    string
	}{
		{
			desc:        "basic",
			given:       jsonresp.InternalServerError(errors.New("nil pointer dereference")),
			wantCode:    http.StatusInternalServerError,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `"Internal Server Error"`,
		},
		{
			desc: "with everything",
			given: jsonresp.InternalServerError(errors.New("database connection failed")).
				WithHeader("X-Request-ID", "req-123456").
				WithCookie(cookie),
			wantCode: http.StatusInternalServerError,
			wantHeaders: map[string]string{
				"X-Request-ID": "req-123456",
			},
			wantCookies: []*http.Cookie{cookie},
			wantBody:    `"Internal Server Error"`,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test-internal-error", nil)

			// When:
			tc.given.Respond(w, r)

			// Then:
			gotCode := w.Code
			if gotCode != tc.wantCode {
				t.Errorf("status code: want %d, got %d", tc.wantCode, gotCode)
			}

			gotHeaders := w.Header()

			// Always check for Content-Type header
			gotContentType := gotHeaders.Get("Content-Type")
			if gotContentType != "application/json" {
				t.Errorf("Content-Type header: want 'application/json', got '%s'", gotContentType)
			}

			// Check other expected headers
			for key, wantValue := range tc.wantHeaders {
				gotValue := gotHeaders.Get(key)
				if gotValue != wantValue {
					t.Errorf("header %s: want %s, got %s", key, wantValue, gotValue)
				}
			}

			gotCookies := w.Result().Cookies()
			if len(gotCookies) != len(tc.wantCookies) {
				t.Errorf("cookie count: want %d, got %d", len(tc.wantCookies), len(gotCookies))
			}

			for i, want := range tc.wantCookies {
				got := gotCookies[i]
				if got.Name != want.Name || got.Value != want.Value {
					t.Errorf("cookie %d: want {%s=%s}, got {%s=%s}",
						i, want.Name, want.Value, got.Name, got.Value)
				}
			}

			gotBody := strings.TrimSpace(w.Body.String())
			if gotBody != tc.wantBody {
				t.Errorf("body: want '%s', got '%s'", tc.wantBody, gotBody)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}

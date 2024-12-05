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
			desc:        "basic",
			given:       jsonresp.Error(errors.New("invalid id"), "Invalid ID provided", http.StatusBadRequest),
			wantCode:    http.StatusBadRequest,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `{"error":"Invalid ID provided"}`,
		},
		{
			desc: "with everything",
			given: jsonresp.Error(errors.New("post not found"), "Post not found", http.StatusNotFound).
				WithHeader("X-Test-1", "test value 1").
				WithCookie(cookie),
			wantCode: http.StatusNotFound,
			wantHeaders: map[string]string{
				"X-Test-1": "test value 1",
			},
			wantCookies: []*http.Cookie{cookie},
			wantBody:    `{"error":"Post not found"}`,
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
			for key, wantValue := range tc.wantHeaders {
				gotValue := gotHeaders.Get(key)
				if gotValue != wantValue {
					t.Errorf("headers %s: want %s, got %s", key, wantValue, gotValue)
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
			wantBody:    `{"error":"Internal Server Error"}`,
		},
		{
			desc: "with eveything",
			given: jsonresp.InternalServerError(errors.New("database failure")).
				WithHeader("X-Test-1", "test value 1").
				WithCookie(cookie),
			wantCode: http.StatusInternalServerError,
			wantHeaders: map[string]string{
				"X-Test-1": "test value 1",
			},
			wantCookies: []*http.Cookie{cookie},
			wantBody:    `{"error":"Internal Server Error"}`,
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
			for key, wantValue := range tc.wantHeaders {
				gotValue := gotHeaders.Get(key)
				if gotValue != wantValue {
					t.Errorf("headers %s: want %s, got %s", key, wantValue, gotValue)
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

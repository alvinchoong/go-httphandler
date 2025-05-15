package jsonresp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/jsonresp"
)

func TestSuccess_Respond(t *testing.T) {
	t.Parallel()

	type User struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at,omitempty"`
	}

	cookie := &http.Cookie{
		Name:  "test-cookie",
		Value: "test-cookie-value",
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
			given:       jsonresp.Success[any](nil),
			wantCode:    http.StatusOK,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `null`,
		},
		{
			desc:        "basic | string",
			given:       jsonresp.Success(ptr("Operation completed successfully")),
			wantCode:    http.StatusOK,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `"Operation completed successfully"`,
		},
		{
			desc:        "basic | map",
			given:       jsonresp.Success(&map[string]string{"message": "Resource deleted successfully"}),
			wantCode:    http.StatusOK,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `{"message":"Resource deleted successfully"}`,
		},
		{
			desc:        "basic | struct",
			given:       jsonresp.Success(&User{ID: "usr_123", Username: "johndoe", Email: "john@example.com", CreatedAt: "2025-01-15T08:30:00Z"}),
			wantCode:    http.StatusOK,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    `{"id":"usr_123","username":"johndoe","email":"john@example.com","created_at":"2025-01-15T08:30:00Z"}`,
		},
		{
			desc: "with everything",
			given: jsonresp.Success(&User{ID: "usr_789", Username: "johndoe", Email: "john@example.com", CreatedAt: "2025-01-15T08:30:00Z"}).
				WithHeader("X-Request-ID", "req-abc123").
				WithHeader("Location", "/api/users/usr_789").
				WithStatus(http.StatusCreated).
				WithCookie(cookie),
			wantCode: http.StatusCreated,
			wantHeaders: map[string]string{
				"X-Request-ID": "req-abc123",
				"Location":     "/api/users/usr_789",
			},
			wantCookies: []*http.Cookie{cookie},
			wantBody:    `{"id":"usr_789","username":"johndoe","email":"john@example.com","created_at":"2025-01-15T08:30:00Z"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/test-success", nil)

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

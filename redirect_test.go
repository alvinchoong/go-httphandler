package httphandler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alvinchoong/go-httphandler"
)

func TestRedirect_Respond(t *testing.T) {
	t.Parallel()

	cookie := &http.Cookie{
		Name:  "test-cookie",
		Value: "test-value",
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
			desc:     "basic",
			given:    httphandler.Redirect("https://google.com", http.StatusMovedPermanently),
			wantCode: http.StatusMovedPermanently,
			wantHeaders: map[string]string{
				"Location": "https://google.com",
			},
			wantCookies: nil,
			wantBody:    `<a href="https://google.com">Moved Permanently</a>.`,
		},
		{
			desc: "with everything",
			given: httphandler.Redirect("https://google.com", http.StatusFound).
				WithHeader("X-Test-1", "test value 1").
				WithCookie(cookie),
			wantCode: http.StatusFound,
			wantHeaders: map[string]string{
				"Location": "https://google.com",
				"X-Test-1": "test value 1",
			},
			wantCookies: []*http.Cookie{cookie},
			wantBody:    `<a href="https://google.com">Found</a>.`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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

package plainresp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/plainresp"
)

func TestSuccess_Respond(t *testing.T) {
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
			given:       plainresp.Success("Success"),
			wantCode:    http.StatusOK,
			wantHeaders: nil,
			wantCookies: nil,
			wantBody:    "Success",
		},
		{
			desc: "with extra headers | with status",
			given: plainresp.Success("OK").
				WithHeader("X-Test-1", "test value 1").
				WithStatus(http.StatusCreated).
				WithCookie(cookie),
			wantCode: http.StatusCreated,
			wantHeaders: map[string]string{
				"X-Test-1": "test value 1",
			},
			wantCookies: []*http.Cookie{cookie},
			wantBody:    "OK",
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

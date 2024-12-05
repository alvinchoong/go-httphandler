package downloadresp_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/downloadresp"
)

func TestInline_Respond(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc        string
		given       httphandler.Responder
		wantCode    int
		wantHeaders map[string]string
		wantBody    string
	}{
		{
			desc:     "txt | without content-type",
			given:    downloadresp.Inline(strings.NewReader("Success"), "test.txt"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/plain; charset=utf-8`,
				"Content-Disposition": `inline; filename="test.txt"`,
			},
			wantBody: "Success",
		},
		{
			desc: "txt | with content-type",
			given: downloadresp.Inline(strings.NewReader("OK"), "test.txt").
				WithContentType("text/plain"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/plain`,
				"Content-Disposition": `inline; filename="test.txt"`,
			},
			wantBody: "OK",
		},
		{
			desc:     "csv | without content-type",
			given:    downloadresp.Inline(strings.NewReader("col1,col2"), "test.csv"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/csv; charset=utf-8`,
				"Content-Disposition": `inline; filename="test.csv"`,
			},
			wantBody: "col1,col2",
		},
		{
			desc: "csv | with content-type",
			given: downloadresp.Inline(strings.NewReader("a,b"), "test.csv").
				WithContentType("text/csv"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/csv`,
				"Content-Disposition": `inline; filename="test.csv"`,
			},
			wantBody: "a,b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			tc.given.Respond(w, r)

			fmt.Printf("%+v\n", w.Header())
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

			gotBody := strings.TrimSpace(w.Body.String())
			if gotBody != tc.wantBody {
				t.Errorf("body: want '%s', got '%s'", tc.wantBody, gotBody)
			}
		})
	}
}

func TestAttachment_Respond(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc        string
		given       httphandler.Responder
		wantCode    int
		wantHeaders map[string]string
		wantBody    string
	}{
		{
			desc:     "txt | without content-type",
			given:    downloadresp.Attachment(strings.NewReader("Success"), "test.txt"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/plain; charset=utf-8`,
				"Content-Disposition": `attachment; filename="test.txt"`,
			},
			wantBody: "Success",
		},
		{
			desc: "txt | with content-type",
			given: downloadresp.Attachment(strings.NewReader("OK"), "test.txt").
				WithContentType("text/plain"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/plain`,
				"Content-Disposition": `attachment; filename="test.txt"`,
			},
			wantBody: "OK",
		},
		{
			desc:     "csv | without content-type",
			given:    downloadresp.Attachment(strings.NewReader("col1,col2"), "test.csv"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/csv; charset=utf-8`,
				"Content-Disposition": `attachment; filename="test.csv"`,
			},
			wantBody: "col1,col2",
		},
		{
			desc: "csv | with content-type",
			given: downloadresp.Attachment(strings.NewReader("a,b"), "test.csv").
				WithContentType("text/csv"),
			wantCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type":        `text/csv`,
				"Content-Disposition": `attachment; filename="test.csv"`,
			},
			wantBody: "a,b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			tc.given.Respond(w, r)

			fmt.Printf("%+v\n", w.Header())
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

			gotBody := strings.TrimSpace(w.Body.String())
			if gotBody != tc.wantBody {
				t.Errorf("body: want '%s', got '%s'", tc.wantBody, gotBody)
			}
		})
	}
}

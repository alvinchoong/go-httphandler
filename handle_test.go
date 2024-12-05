package httphandler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvinchoong/go-httphandler"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc     string
		given    httphandler.RequestHandler
		wantCode int
		wantBody string
	}{
		{
			desc: "handle success",
			given: func(r *http.Request) httphandler.Responder {
				return &mockResponder{
					StatusCode: http.StatusOK,
					Body:       "Success",
				}
			},
			wantCode: http.StatusOK,
			wantBody: "Success",
		},
		{
			desc: "handle nil",
			given: func(r *http.Request) httphandler.Responder {
				return nil
			},
			wantCode: http.StatusNoContent,
			wantBody: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			// Given:
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			// When:
			httphandler.Handle(tc.given).ServeHTTP(w, r)

			// Then:
			gotCode := w.Code
			if gotCode != tc.wantCode {
				t.Errorf("status code: want %d, got %d", tc.wantCode, gotCode)
			}

			gotBody := w.Body.String()
			if gotBody != tc.wantBody {
				t.Errorf("body: want '%s', got '%s'", tc.wantBody, gotBody)
			}
		})
	}
}

func TestHandleWithInput(t *testing.T) {
	type Input struct {
		Name string
		Age  int
	}

	testCases := []struct {
		name           string
		decode         httphandler.RequestDecodeFunc[Input]
		handler        httphandler.RequestHandlerWithInput[Input]
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "decode success | handle success",
			decode: func(r *http.Request) (Input, error) {
				return Input{Name: "Alice", Age: 30}, nil
			},
			handler: func(r *http.Request, got Input) httphandler.Responder {
				want := Input{
					Name: "Alice",
					Age:  30,
				}
				if got != want {
					t.Errorf("handler input: want %+v, got %+v", want, got)
				}
				return &mockResponder{
					StatusCode: http.StatusOK,
					Body:       "Success",
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Success",
		},
		{
			name: "decode fail",
			decode: func(r *http.Request) (Input, error) {
				return Input{}, errors.New("decoding failed")
			},
			handler: func(r *http.Request, input Input) httphandler.Responder {
				t.Errorf("handler: should not be called on decoding failure")
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request payload\n",
		},
		{
			name: "decode success | handler nil",
			decode: func(r *http.Request) (Input, error) {
				return Input{Name: "Bob", Age: 25}, nil
			},
			handler: func(r *http.Request, got Input) httphandler.Responder {
				want := Input{
					Name: "Bob",
					Age:  25,
				}
				if got != want {
					t.Errorf("handler input: want %+v, got %+v", want, got)
				}
				return nil
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Given:
			r := httptest.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()

			given := httphandler.HandleWithInput(tc.handler, httphandler.WithDecodeFunc(tc.decode))

			// When:
			given.ServeHTTP(w, r)

			// Then:
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tc.expectedBody, w.Body.String())
			}
		})
	}
}

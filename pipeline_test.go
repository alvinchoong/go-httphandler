package httphandler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPipelineCompilation validates that the pipeline types compile correctly
func TestPipelineCompilation(t *testing.T) {
	// Define test types
	type TestContext1 struct{ Value string }
	type TestContext2 struct{ Value string }
	type TestContext3 struct{ Value string }
	type TestContext4 struct{ Value string }
	type TestContext5 struct{ Value string }
	type TestContext6 struct{ Value string }
	type TestContext7 struct{ Value string }
	type TestContext8 struct{ Value string }
	type TestInput struct{ Value string }

	// Define stub decoders
	decoder1 := func(r *http.Request) (TestContext1, error) {
		return TestContext1{Value: "ctx1"}, nil
	}

	decoder2 := func(r *http.Request, ctx1 TestContext1) (TestContext2, error) {
		return TestContext2{Value: "ctx2"}, nil
	}

	decoder3 := func(r *http.Request, ctx1 TestContext1, ctx2 TestContext2) (TestContext3, error) {
		return TestContext3{Value: "ctx3"}, nil
	}

	decoder4 := func(r *http.Request, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3) (TestContext4, error) {
		return TestContext4{Value: "ctx4"}, nil
	}

	decoder5 := func(r *http.Request, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4) (TestContext5, error) {
		return TestContext5{Value: "ctx5"}, nil
	}

	decoder6 := func(r *http.Request, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5) (TestContext6, error) {
		return TestContext6{Value: "ctx6"}, nil
	}

	decoder7 := func(r *http.Request, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6) (TestContext7, error) {
		return TestContext7{Value: "ctx7"}, nil
	}

	decoder8 := func(r *http.Request, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6, ctx7 TestContext7) (TestContext8, error) {
		return TestContext8{Value: "ctx8"}, nil
	}

	inputDecoder := func(r *http.Request) (TestInput, error) {
		return TestInput{Value: "input"}, nil
	}

	// Create test pipeline chains
	p1 := NewPipeline1(decoder1)
	p2 := NewPipeline2(p1, decoder2)
	p3 := NewPipeline3(p2, decoder3)
	p4 := NewPipeline4(p3, decoder4)
	p5 := NewPipeline5(p4, decoder5)
	p6 := NewPipeline6(p5, decoder6)
	p7 := NewPipeline7(p6, decoder7)
	p8 := NewPipeline8(p7, decoder8)

	// Create handlers (just verifying compilation)
	_ = HandlePipelineWithInput1(p1, inputDecoder, func(ctx1 TestContext1, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput2(p2, inputDecoder, func(ctx1 TestContext1, ctx2 TestContext2, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput3(p3, inputDecoder, func(ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput4(p4, inputDecoder, func(ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput5(p5, inputDecoder, func(ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput6(p6, inputDecoder, func(ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput7(p7, inputDecoder, func(ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6, ctx7 TestContext7, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput8(p8, inputDecoder, func(ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6, ctx7 TestContext7, ctx8 TestContext8, input TestInput) Responder {
		return nil
	})

	// If this test compiles, it passes
}

// TestPipelineExecution tests the actual execution of pipelines
func TestPipelineExecution(t *testing.T) {
	// Test context and input types
	type UserContext struct {
		Username string
	}

	type ActionContext struct {
		Action string
	}

	type LoginInput struct {
		Password string
	}

	// Test decoders
	decodeUser := func(r *http.Request) (UserContext, error) {
		// Extract username from Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			return UserContext{}, errors.New("missing authorization header")
		}
		return UserContext{Username: auth}, nil
	}

	decodeAction := func(r *http.Request, user UserContext) (ActionContext, error) {
		// Extract action from URL path
		action := r.URL.Path
		if action == "" {
			return ActionContext{}, errors.New("missing action")
		}
		return ActionContext{Action: action}, nil
	}

	decodeLoginInput := func(r *http.Request) (LoginInput, error) {
		// In a real implementation, we'd parse JSON or form data
		// For this test, we're just using a simple string
		return LoginInput{Password: "test-password"}, nil
	}

	// Create pipeline
	userPipeline := NewPipeline1(decodeUser)
	actionPipeline := NewPipeline2(userPipeline, decodeAction)

	// Create handler
	handler := HandlePipelineWithInput2(actionPipeline, decodeLoginInput,
		func(user UserContext, action ActionContext, input LoginInput) Responder {
			// Simple success responder for testing
			return &testResponder{
				message: "Success: " + user.Username + " " + action.Action + " " + input.Password,
			}
		})

	// Create test request
	req := httptest.NewRequest("POST", "/login", strings.NewReader(""))
	req.Header.Set("Authorization", "test-user")

	// Create test response recorder
	w := httptest.NewRecorder()

	// Execute handler
	handler(w, req)

	// Check result
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := "Success: test-user /login test-password"
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
	}
}

// TestPipelineErrorHandling tests error handling in pipelines
func TestPipelineErrorHandling(t *testing.T) {
	// Test decoder that always fails
	failingDecoder := func(r *http.Request) (struct{}, error) {
		return struct{}{}, errors.New("decoder error")
	}

	// Create pipeline with failing decoder
	pipeline := NewPipeline1(failingDecoder)

	// Create handler
	handler := HandlePipelineWithInput1(pipeline, func(r *http.Request) (struct{}, error) {
		return struct{}{}, nil
	}, func(ctx struct{}, input struct{}) Responder {
		return &testResponder{message: "This should not be called"}
	})

	// Create test request
	req := httptest.NewRequest("GET", "/", nil)

	// Create test response recorder
	w := httptest.NewRecorder()

	// Execute handler
	handler(w, req)

	// Check error handling
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Should contain the error message
	if !strings.Contains(w.Body.String(), "decoder error") {
		t.Errorf("expected error message to contain 'decoder error', got %q", w.Body.String())
	}
}

// testResponder is a simple Responder implementation for testing
type testResponder struct {
	message string
}

func (r *testResponder) Respond(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.message))
}

// TestCustomErrorHandler tests the custom error handling capabilities
func TestCustomErrorHandler(t *testing.T) {
	// Custom error handler for context errors
	custom401Handler := func(stage int, err error) Responder {
		return &customErrorResponder{
			statusCode: http.StatusUnauthorized,
			message:   fmt.Sprintf("Auth failed at stage %d: %v", stage, err),
		}
	}

	// Custom error handler for input errors
	custom422Handler := func(err error) Responder {
		return &customErrorResponder{
			statusCode: http.StatusUnprocessableEntity,
			message:   fmt.Sprintf("Invalid input: %v", err),
		}
	}

	// Test decoder that always fails
	failingDecoder := func(r *http.Request) (struct{}, error) {
		return struct{}{}, errors.New("auth token invalid")
	}

	// Test input decoder that always fails
	failingInputDecoder := func(r *http.Request) (struct{}, error) {
		return struct{}{}, errors.New("missing required field")
	}

	// Test with custom context error handler
	pipeline1 := NewPipeline1(
		failingDecoder,
		WithContextErrorHandler(custom401Handler),
	)

	handler1 := HandlePipelineWithInput1(pipeline1, func(r *http.Request) (struct{}, error) {
		return struct{}{}, nil
	}, func(ctx struct{}, input struct{}) Responder {
		return &testResponder{message: "This should not be called"}
	})

	// Test with custom input error handler
	pipeline2 := NewPipeline1(
		func(r *http.Request) (struct{}, error) { return struct{}{}, nil },
		WithInputErrorHandler(custom422Handler),
	)

	handler2 := HandlePipelineWithInput1(pipeline2, failingInputDecoder, func(ctx struct{}, input struct{}) Responder {
		return &testResponder{message: "This should not be called"}
	})

	// Test context error handler
	t.Run("CustomContextErrorHandler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler1(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		}

		expected := "Auth failed at stage 1: auth token invalid"
		if !strings.Contains(w.Body.String(), expected) {
			t.Errorf("expected body to contain '%s', got %q", expected, w.Body.String())
		}
	})

	// Test input error handler
	t.Run("CustomInputErrorHandler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler2(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
		}

		expected := "Invalid input: missing required field"
		if !strings.Contains(w.Body.String(), expected) {
			t.Errorf("expected body to contain '%s', got %q", expected, w.Body.String())
		}
	})
}

// Custom error responder for testing
type customErrorResponder struct {
	statusCode int
	message    string
}

func (r *customErrorResponder) Respond(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(r.statusCode)
	w.Write([]byte(r.message))
}

// TestPipelineDeepChaining tests the full pipeline depth
func TestPipelineDeepChaining(t *testing.T) {
	// Define test types for all 8 contexts
	type Context1 struct{ Value string }
	type Context2 struct{ Value string }
	type Context3 struct{ Value string }
	type Context4 struct{ Value string }
	type Context5 struct{ Value string }
	type Context6 struct{ Value string }
	type Context7 struct{ Value string }
	type Context8 struct{ Value string }
	type Input struct{ Value string }

	// Context decoders
	decoder1 := func(r *http.Request) (Context1, error) {
		return Context1{Value: "ctx1"}, nil
	}
	decoder2 := func(r *http.Request, c1 Context1) (Context2, error) {
		return Context2{Value: c1.Value + "-ctx2"}, nil
	}
	decoder3 := func(r *http.Request, c1 Context1, c2 Context2) (Context3, error) {
		return Context3{Value: c2.Value + "-ctx3"}, nil
	}
	decoder4 := func(r *http.Request, c1 Context1, c2 Context2, c3 Context3) (Context4, error) {
		return Context4{Value: c3.Value + "-ctx4"}, nil
	}
	decoder5 := func(r *http.Request, c1 Context1, c2 Context2, c3 Context3, c4 Context4) (Context5, error) {
		return Context5{Value: c4.Value + "-ctx5"}, nil
	}
	decoder6 := func(r *http.Request, c1 Context1, c2 Context2, c3 Context3, c4 Context4, c5 Context5) (Context6, error) {
		return Context6{Value: c5.Value + "-ctx6"}, nil
	}
	decoder7 := func(r *http.Request, c1 Context1, c2 Context2, c3 Context3, c4 Context4, c5 Context5, c6 Context6) (Context7, error) {
		return Context7{Value: c6.Value + "-ctx7"}, nil
	}
	decoder8 := func(r *http.Request, c1 Context1, c2 Context2, c3 Context3, c4 Context4, c5 Context5, c6 Context6, c7 Context7) (Context8, error) {
		return Context8{Value: c7.Value + "-ctx8"}, nil
	}

	// Input decoder
	inputDecoder := func(r *http.Request) (Input, error) {
		return Input{Value: "input"}, nil
	}

	// Create the deepest pipeline
	p1 := NewPipeline1(decoder1)
	p2 := NewPipeline2(p1, decoder2)
	p3 := NewPipeline3(p2, decoder3)
	p4 := NewPipeline4(p3, decoder4)
	p5 := NewPipeline5(p4, decoder5)
	p6 := NewPipeline6(p5, decoder6)
	p7 := NewPipeline7(p6, decoder7)
	p8 := NewPipeline8(p7, decoder8)

	// Create handler with all 8 contexts
	handler := HandlePipelineWithInput8(p8, inputDecoder, 
		func(c1 Context1, c2 Context2, c3 Context3, c4 Context4, c5 Context5, c6 Context6, c7 Context7, c8 Context8, input Input) Responder {
			return &testResponder{
				message: fmt.Sprintf("c1=%s, c2=%s, c3=%s, c4=%s, c5=%s, c6=%s, c7=%s, c8=%s, input=%s",
					c1.Value, c2.Value, c3.Value, c4.Value, c5.Value, c6.Value, c7.Value, c8.Value, input.Value),
			}
		},
	)

	// Test the chained pipeline
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expected := "c1=ctx1, c2=ctx1-ctx2, c3=ctx1-ctx2-ctx3, c4=ctx1-ctx2-ctx3-ctx4, c5=ctx1-ctx2-ctx3-ctx4-ctx5, c6=ctx1-ctx2-ctx3-ctx4-ctx5-ctx6, c7=ctx1-ctx2-ctx3-ctx4-ctx5-ctx6-ctx7, c8=ctx1-ctx2-ctx3-ctx4-ctx5-ctx6-ctx7-ctx8, input=input"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}

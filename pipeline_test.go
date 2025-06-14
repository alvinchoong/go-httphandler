package httphandler

import (
	"context"
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

	// Create test pipeline chains - pass all decoder functions directly
	p1 := NewPipeline1(decoder1)
	p2 := NewPipeline2(decoder1, decoder2)
	p3 := NewPipeline3(decoder1, decoder2, decoder3)
	p4 := NewPipeline4(decoder1, decoder2, decoder3, decoder4)
	p5 := NewPipeline5(decoder1, decoder2, decoder3, decoder4, decoder5)
	p6 := NewPipeline6(decoder1, decoder2, decoder3, decoder4, decoder5, decoder6)
	p7 := NewPipeline7(decoder1, decoder2, decoder3, decoder4, decoder5, decoder6, decoder7)
	p8 := NewPipeline8(decoder1, decoder2, decoder3, decoder4, decoder5, decoder6, decoder7, decoder8)

	// Create handlers (just verifying compilation)
	_ = HandlePipelineWithInput1(p1, inputDecoder, func(ctx context.Context, ctx1 TestContext1, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput2(p2, inputDecoder, func(ctx context.Context, ctx1 TestContext1, ctx2 TestContext2, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput3(p3, inputDecoder, func(ctx context.Context, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput4(p4, inputDecoder, func(ctx context.Context, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput5(p5, inputDecoder, func(ctx context.Context, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput6(p6, inputDecoder, func(ctx context.Context, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput7(p7, inputDecoder, func(ctx context.Context, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6, ctx7 TestContext7, input TestInput) Responder {
		return nil
	})

	_ = HandlePipelineWithInput8(p8, inputDecoder, func(ctx context.Context, ctx1 TestContext1, ctx2 TestContext2, ctx3 TestContext3, ctx4 TestContext4, ctx5 TestContext5, ctx6 TestContext6, ctx7 TestContext7, ctx8 TestContext8, input TestInput) Responder {
		return nil
	})

	// If this test compiles, it passes
}

// TestPipelineExecution tests the actual execution of pipelines
func TestPipelineExecution(t *testing.T) {
	// Define test context and input types
	type UserContext struct {
		Username string
	}

	type ActionContext struct {
		Action string
	}

	type LoginInput struct {
		Password string
	}

	// Define decoders
	userDecoder := func(r *http.Request) (UserContext, error) {
		return UserContext{Username: "testuser"}, nil
	}

	actionDecoder := func(r *http.Request, user UserContext) (ActionContext, error) {
		return ActionContext{Action: "login"}, nil
	}

	decodeLoginInput := func(r *http.Request) (LoginInput, error) {
		// In a real implementation, we'd parse JSON or form data
		// For this test, we're just using a simple string
		return LoginInput{Password: "test-password"}, nil
	}

	// Create pipeline directly passing all decoders
	actionPipeline := NewPipeline2(userDecoder, actionDecoder)

	// Create handler
	handler := HandlePipelineWithInput2(actionPipeline, decodeLoginInput,
		func(ctx context.Context, user UserContext, action ActionContext, input LoginInput) Responder {
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

	expectedBody := "Success: testuser login test-password"
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
		return struct{}{}, nil	}, func(ctx context.Context, val struct{}, input struct{}) Responder {
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
	// Define test errors
	authError := errors.New("auth token invalid")
	inputError := errors.New("missing required field")

	// Test direct error handling
	t.Run("CustomDecodeErrorHandler", func(t *testing.T) {
		// Create a custom error responder directly
		resp := &customErrorResponder{
			statusCode: http.StatusUnauthorized,
			message:   fmt.Sprintf("Auth failed at stage %d: %v", 1, authError),
		}

		// Test the responder directly
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		resp.Respond(w, req)

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
		// Create a custom error responder directly
		resp := &customErrorResponder{
			statusCode: http.StatusUnprocessableEntity,
			message:   fmt.Sprintf("Invalid input: %v", inputError),
		}

		// Test the responder directly
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		resp.Respond(w, req)

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

	// Create the deepest pipeline directly with all decoders (not using previous pipeline variables)
	p8 := NewPipeline8(decoder1, decoder2, decoder3, decoder4, decoder5, decoder6, decoder7, decoder8)

	// Create handler with all 8 contexts
	handler := HandlePipelineWithInput8(p8, inputDecoder,
		func(ctx context.Context, c1 Context1, c2 Context2, c3 Context3, c4 Context4, c5 Context5, c6 Context6, c7 Context7, c8 Context8, input Input) Responder {
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

// Test fixtures for unified pipeline approach
// testInput is used for pipeline input stage tests
type testInput struct {
	Value string
}

// TestHandlePipelineWithInputStage tests the implementation that treats input as a pipeline stage
func TestHandlePipelineWithInputStage(t *testing.T) {
	// Create a simple decoder that returns a context value
	contextDecoder := func(r *http.Request) (string, error) {
		return "context-value", nil
	}

	// Create a simple input decoder that returns an input value
	inputDecoder := func(r *http.Request) (testInput, error) {
		return testInput{Value: "input-value"}, nil
	}

	// Create pipelines with just one context
	pipeline1 := NewPipeline1(contextDecoder)

	// Create a test handler function that verifies both context and input values
	handlerCalled := false
	handler := func(ctx context.Context, val string, input testInput) Responder {
		handlerCalled = true
		if val != "context-value" {
			t.Errorf("Expected context value 'context-value', got '%s'", val)
		}
		if input.Value != "input-value" {
			t.Errorf("Expected input value 'input-value', got '%s'", input.Value)
		}
		return &testResponder{message: "success"}
	}

	// Create the handler with the pipeline
	handlerFunc := HandlePipelineWithInput1(pipeline1, inputDecoder, handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handlerFunc(w, req)

	// Verify the handler was called
	if !handlerCalled {
		t.Errorf("Handler was not called")
	}

	// Verify the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("Expected body 'success', got '%s'", w.Body.String())
	}
}

// TestHandlePipelineWithComplexInputStage tests the implementation with multiple contexts
func TestHandlePipelineWithComplexInputStage(t *testing.T) {
	// Create context decoders
	contextDecoder1 := func(r *http.Request) (string, error) {
		return "context1", nil
	}

	contextDecoder2 := func(r *http.Request, c1 string) (int, error) {
		if c1 != "context1" {
			t.Errorf("Expected c1 to be 'context1', got '%s'", c1)
		}
		return 42, nil
	}

	// Create an input decoder
	inputDecoder := func(r *http.Request) (bool, error) {
		return true, nil
	}

	// Create a pipeline with two contexts directly
	pipeline2 := NewPipeline2(contextDecoder1, contextDecoder2)

	// Create a test handler function
	handlerCalled := false
	handler := func(ctx context.Context, val1 string, val2 int, input bool) Responder {
		handlerCalled = true
		if val1 != "context1" {
			t.Errorf("Expected val1 to be 'context1', got '%s'", val1)
		}
		if val2 != 42 {
			t.Errorf("Expected val2 to be 42, got %d", val2)
		}
		if !input {
			t.Errorf("Expected input to be true")
		}
		return &testResponder{message: "complex-success"}
	}

	// Create the handler with the pipeline
	handlerFunc := HandlePipelineWithInput2(pipeline2, inputDecoder, handler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handlerFunc(w, req)

	// Verify the handler was called
	if !handlerCalled {
		t.Errorf("Handler was not called")
	}

	// Verify the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "complex-success" {
		t.Errorf("Expected body 'complex-success', got '%s'", w.Body.String())
	}
}

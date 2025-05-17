package httphandler

import (
	"context"
	"fmt"
	"net/http"
)

// ========== Regular pipeline handlers ==========

// Default error handlers when options is nil
func handleContextError(options *PipelineOptions, stage int, err error) Responder {
	if options == nil || options.ContextErrorHandler == nil {
		return defaultErrorHandler(err)
	}
	return options.ContextErrorHandler(stage, err)
}

func handleInputError(options *PipelineOptions, err error) Responder {
	if options == nil || options.InputErrorHandler == nil {
		return defaultErrorHandler(err)
	}
	return options.InputErrorHandler(err)
}

// ResponderFunc is a function type that implements the Responder interface
type ResponderFunc func(w http.ResponseWriter, r *http.Request)

// Respond implements the Responder interface
func (f ResponderFunc) Respond(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

// Default error handler returns a 400 Bad Request with the error message
func defaultErrorHandler(err error) Responder {
	return ResponderFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
	})
}

// HandlePipeline1 creates a handler using a pipeline with one context
func HandlePipeline1[C any](
	p Pipeline1[C],
	handler func(ctx context.Context, val C) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode context
		val, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipeline2 creates a handler using a pipeline with two contexts
func HandlePipeline2[C1, C2 any](
	p Pipeline2[C1, C2],
	handler func(ctx context.Context, val1 C1, val2 C2) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.decoder2(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipeline3 creates a handler using a pipeline with three contexts
func HandlePipeline3[C1, C2, C3 any](
	p Pipeline3[C1, C2, C3],
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.decoder2(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := p.decoder3(r, val1, val2)
		if err != nil {
			p.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipeline4 creates a handler using a pipeline with four contexts
func HandlePipeline4[C1, C2, C3, C4 any](
	p Pipeline4[C1, C2, C3, C4],
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.decoder2(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := p.decoder3(r, val1, val2)
		if err != nil {
			p.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := p.decoder4(r, val1, val2, val3)
		if err != nil {
			p.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipeline5 creates a handler using a pipeline with five contexts
func HandlePipeline5[C1, C2, C3, C4, C5 any](
	p Pipeline5[C1, C2, C3, C4, C5],
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.decoder2(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := p.decoder3(r, val1, val2)
		if err != nil {
			p.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := p.decoder4(r, val1, val2, val3)
		if err != nil {
			p.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := p.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			p.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipeline6 creates a handler using a pipeline with six contexts
func HandlePipeline6[C1, C2, C3, C4, C5, C6 any](
	p Pipeline6[C1, C2, C3, C4, C5, C6],
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.decoder2(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := p.decoder3(r, val1, val2)
		if err != nil {
			p.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := p.decoder4(r, val1, val2, val3)
		if err != nil {
			p.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := p.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			p.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := p.decoder6(r, val1, val2, val3, val4, val5)
		if err != nil {
			p.options.ContextErrorHandler(6, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5, val6)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipeline7 creates a handler using a pipeline with seven contexts
func HandlePipeline7[C1, C2, C3, C4, C5, C6, C7 any](
	p Pipeline7[C1, C2, C3, C4, C5, C6, C7],
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.decoder2(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := p.decoder3(r, val1, val2)
		if err != nil {
			p.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := p.decoder4(r, val1, val2, val3)
		if err != nil {
			p.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := p.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			p.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := p.decoder6(r, val1, val2, val3, val4, val5)
		if err != nil {
			p.options.ContextErrorHandler(6, err).Respond(w, r)
			return
		}

		// Decode seventh context
		val7, err := p.decoder7(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			p.options.ContextErrorHandler(7, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5, val6, val7)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipeline8 creates a handler using a pipeline with eight contexts
func HandlePipeline8[C1, C2, C3, C4, C5, C6, C7, C8 any](
	p Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8],
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7, val8 C8) Responder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.decoder1(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.decoder2(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := p.decoder3(r, val1, val2)
		if err != nil {
			p.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := p.decoder4(r, val1, val2, val3)
		if err != nil {
			p.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := p.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			p.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := p.decoder6(r, val1, val2, val3, val4, val5)
		if err != nil {
			p.options.ContextErrorHandler(6, err).Respond(w, r)
			return
		}

		// Decode seventh context
		val7, err := p.decoder7(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			p.options.ContextErrorHandler(7, err).Respond(w, r)
			return
		}

		// Decode eighth context
		val8, err := p.decoder8(r, val1, val2, val3, val4, val5, val6, val7)
		if err != nil {
			p.options.ContextErrorHandler(8, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5, val6, val7, val8)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// ========== Input as pipeline stage functions ==========

// NewPipelineWithInput1 creates a pipeline with one context type and input
func NewPipelineWithInput1[C, T any](
	p Pipeline1[C],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput1[C, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}
	
	return PipelineWithInput1[C, T]{
		decoder1: p.decoder1,
		decoder2: func(r *http.Request, c C) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// NewPipelineWithInput2 creates a pipeline with two context types and input
func NewPipelineWithInput2[C1, C2, T any](
	p Pipeline2[C1, C2],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput2[C1, C2, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}

	return PipelineWithInput2[C1, C2, T]{
		decoder1: p.decoder1,
		decoder2: p.decoder2,
		decoder3: func(r *http.Request, c1 C1, c2 C2) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// NewPipelineWithInput3 creates a pipeline with three context types and input
func NewPipelineWithInput3[C1, C2, C3, T any](
	p Pipeline3[C1, C2, C3],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput3[C1, C2, C3, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}

	return PipelineWithInput3[C1, C2, C3, T]{
		decoder1: p.decoder1,
		decoder2: p.decoder2,
		decoder3: p.decoder3,
		decoder4: func(r *http.Request, c1 C1, c2 C2, c3 C3) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// NewPipelineWithInput4 creates a pipeline with four context types and input
func NewPipelineWithInput4[C1, C2, C3, C4, T any](
	p Pipeline4[C1, C2, C3, C4],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput4[C1, C2, C3, C4, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}

	return PipelineWithInput4[C1, C2, C3, C4, T]{
		decoder1: p.decoder1,
		decoder2: p.decoder2,
		decoder3: p.decoder3,
		decoder4: p.decoder4,
		decoder5: func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// NewPipelineWithInput5 creates a pipeline with five context types and input
func NewPipelineWithInput5[C1, C2, C3, C4, C5, T any](
	p Pipeline5[C1, C2, C3, C4, C5],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput5[C1, C2, C3, C4, C5, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}

	return PipelineWithInput5[C1, C2, C3, C4, C5, T]{
		decoder1: p.decoder1,
		decoder2: p.decoder2,
		decoder3: p.decoder3,
		decoder4: p.decoder4,
		decoder5: p.decoder5,
		decoder6: func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// NewPipelineWithInput6 creates a pipeline with six context types and input
func NewPipelineWithInput6[C1, C2, C3, C4, C5, C6, T any](
	p Pipeline6[C1, C2, C3, C4, C5, C6],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput6[C1, C2, C3, C4, C5, C6, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}

	return PipelineWithInput6[C1, C2, C3, C4, C5, C6, T]{
		decoder1: p.decoder1,
		decoder2: p.decoder2,
		decoder3: p.decoder3,
		decoder4: p.decoder4,
		decoder5: p.decoder5,
		decoder6: p.decoder6,
		decoder7: func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// NewPipelineWithInput7 creates a pipeline with seven context types and input
func NewPipelineWithInput7[C1, C2, C3, C4, C5, C6, C7, T any](
	p Pipeline7[C1, C2, C3, C4, C5, C6, C7],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput7[C1, C2, C3, C4, C5, C6, C7, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}

	return PipelineWithInput7[C1, C2, C3, C4, C5, C6, C7, T]{
		decoder1: p.decoder1,
		decoder2: p.decoder2,
		decoder3: p.decoder3,
		decoder4: p.decoder4,
		decoder5: p.decoder5,
		decoder6: p.decoder6,
		decoder7: p.decoder7,
		decoder8: func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// NewPipelineWithInput8 creates a pipeline with eight context types and input
func NewPipelineWithInput8[C1, C2, C3, C4, C5, C6, C7, C8, T any](
	p Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) PipelineWithInput8[C1, C2, C3, C4, C5, C6, C7, C8, T] {
	// Apply all provided options to a new options struct
	opts := &PipelineOptions{}
	for _, option := range options {
		option(opts)
	}

	return PipelineWithInput8[C1, C2, C3, C4, C5, C6, C7, C8, T]{
		decoder1: p.decoder1,
		decoder2: p.decoder2,
		decoder3: p.decoder3,
		decoder4: p.decoder4,
		decoder5: p.decoder5,
		decoder6: p.decoder6,
		decoder7: p.decoder7,
		decoder8: p.decoder8,
		decoder9: func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7, c8 C8) (T, error) {
			return inputDecoder(r)
		},
		options: *opts,
	}
}

// PipelineWithInput1 is a pipeline stage with one context and input
type PipelineWithInput1[C, T any] struct {
	decoder1  func(r *http.Request) (C, error)
	decoder2  func(r *http.Request, val1 C) (T, error)
	options   PipelineOptions
}

// PipelineWithInput2 is a pipeline stage with two contexts and input
type PipelineWithInput2[C1, C2, T any] struct {
	decoder1  func(r *http.Request) (C1, error)
	decoder2  func(r *http.Request, val1 C1) (C2, error)
	decoder3  func(r *http.Request, val1 C1, val2 C2) (T, error)
	options   PipelineOptions
}

// PipelineWithInput3 is a pipeline stage with three contexts and input
type PipelineWithInput3[C1, C2, C3, T any] struct {
	decoder1  func(r *http.Request) (C1, error)
	decoder2  func(r *http.Request, val1 C1) (C2, error)
	decoder3  func(r *http.Request, val1 C1, val2 C2) (C3, error)
	decoder4  func(r *http.Request, val1 C1, val2 C2, val3 C3) (T, error)
	options   PipelineOptions
}

// PipelineWithInput4 is a pipeline stage with four contexts and input
type PipelineWithInput4[C1, C2, C3, C4, T any] struct {
	decoder1  func(r *http.Request) (C1, error)
	decoder2  func(r *http.Request, val1 C1) (C2, error)
	decoder3  func(r *http.Request, val1 C1, val2 C2) (C3, error)
	decoder4  func(r *http.Request, val1 C1, val2 C2, val3 C3) (C4, error)
	decoder5  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4) (T, error)
	options   PipelineOptions
}

// PipelineWithInput5 is a pipeline stage with five contexts and input
type PipelineWithInput5[C1, C2, C3, C4, C5, T any] struct {
	decoder1  func(r *http.Request) (C1, error)
	decoder2  func(r *http.Request, val1 C1) (C2, error)
	decoder3  func(r *http.Request, val1 C1, val2 C2) (C3, error)
	decoder4  func(r *http.Request, val1 C1, val2 C2, val3 C3) (C4, error)
	decoder5  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4) (C5, error)
	decoder6  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5) (T, error)
	options   PipelineOptions
}

// PipelineWithInput6 is a pipeline stage with six contexts and input
type PipelineWithInput6[C1, C2, C3, C4, C5, C6, T any] struct {
	decoder1  func(r *http.Request) (C1, error)
	decoder2  func(r *http.Request, val1 C1) (C2, error)
	decoder3  func(r *http.Request, val1 C1, val2 C2) (C3, error)
	decoder4  func(r *http.Request, val1 C1, val2 C2, val3 C3) (C4, error)
	decoder5  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4) (C5, error)
	decoder6  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5) (C6, error)
	decoder7  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6) (T, error)
	options   PipelineOptions
}

// PipelineWithInput7 is a pipeline stage with seven contexts and input
type PipelineWithInput7[C1, C2, C3, C4, C5, C6, C7, T any] struct {
	decoder1  func(r *http.Request) (C1, error)
	decoder2  func(r *http.Request, val1 C1) (C2, error)
	decoder3  func(r *http.Request, val1 C1, val2 C2) (C3, error)
	decoder4  func(r *http.Request, val1 C1, val2 C2, val3 C3) (C4, error)
	decoder5  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4) (C5, error)
	decoder6  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5) (C6, error)
	decoder7  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6) (C7, error)
	decoder8  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7) (T, error)
	options   PipelineOptions
}

// PipelineWithInput8 is a pipeline stage with eight contexts and input
type PipelineWithInput8[C1, C2, C3, C4, C5, C6, C7, C8, T any] struct {
	decoder1  func(r *http.Request) (C1, error)
	decoder2  func(r *http.Request, val1 C1) (C2, error)
	decoder3  func(r *http.Request, val1 C1, val2 C2) (C3, error)
	decoder4  func(r *http.Request, val1 C1, val2 C2, val3 C3) (C4, error)
	decoder5  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4) (C5, error)
	decoder6  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5) (C6, error)
	decoder7  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6) (C7, error)
	decoder8  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7) (C8, error)
	decoder9  func(r *http.Request, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7, val8 C8) (T, error)
	options   PipelineOptions
}

// ========== Handler functions with input as final pipeline stage ==========

// HandlePipelineWithInput1 creates a handler with one context and input as a pipeline stage
func HandlePipelineWithInput1[C, T any](
	p Pipeline1[C],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val C, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput1(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode input (as second context)
		input, err := pipeline.decoder2(r, val1)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipelineWithInput2 creates a handler with two contexts and input as a pipeline stage
func HandlePipelineWithInput2[C1, C2, T any](
	p Pipeline2[C1, C2],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val1 C1, val2 C2, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput2(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.decoder2(r, val1)
		if err != nil {
			handleContextError(&pipeline.options, 2, err).Respond(w, r)
			return
		}

		// Decode input (as third context)
		input, err := pipeline.decoder3(r, val1, val2)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipelineWithInput3 creates a handler with three contexts and input as a pipeline stage
func HandlePipelineWithInput3[C1, C2, C3, T any](
	p Pipeline3[C1, C2, C3],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput3(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.decoder2(r, val1)
		if err != nil {
			handleContextError(&pipeline.options, 2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.decoder3(r, val1, val2)
		if err != nil {
			handleContextError(&pipeline.options, 3, err).Respond(w, r)
			return
		}

		// Decode input (as fourth context)
		input, err := pipeline.decoder4(r, val1, val2, val3)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipelineWithInput4 creates a handler with four contexts and input as a pipeline stage
func HandlePipelineWithInput4[C1, C2, C3, C4, T any](
	p Pipeline4[C1, C2, C3, C4],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput4(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.decoder2(r, val1)
		if err != nil {
			handleContextError(&pipeline.options, 2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.decoder3(r, val1, val2)
		if err != nil {
			handleContextError(&pipeline.options, 3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.decoder4(r, val1, val2, val3)
		if err != nil {
			handleContextError(&pipeline.options, 4, err).Respond(w, r)
			return
		}

		// Decode input (as fifth context)
		input, err := pipeline.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipelineWithInput5 creates a handler with five contexts and input as a pipeline stage
func HandlePipelineWithInput5[C1, C2, C3, C4, C5, T any](
	p Pipeline5[C1, C2, C3, C4, C5],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput5(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.decoder2(r, val1)
		if err != nil {
			handleContextError(&pipeline.options, 2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.decoder3(r, val1, val2)
		if err != nil {
			handleContextError(&pipeline.options, 3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.decoder4(r, val1, val2, val3)
		if err != nil {
			handleContextError(&pipeline.options, 4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := pipeline.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			handleContextError(&pipeline.options, 5, err).Respond(w, r)
			return
		}

		// Decode input (as sixth context)
		input, err := pipeline.decoder6(r, val1, val2, val3, val4, val5)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipelineWithInput6 creates a handler with six contexts and input as a pipeline stage
func HandlePipelineWithInput6[C1, C2, C3, C4, C5, C6, T any](
	p Pipeline6[C1, C2, C3, C4, C5, C6],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput6(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.decoder2(r, val1)
		if err != nil {
			handleContextError(&pipeline.options, 2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.decoder3(r, val1, val2)
		if err != nil {
			handleContextError(&pipeline.options, 3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.decoder4(r, val1, val2, val3)
		if err != nil {
			handleContextError(&pipeline.options, 4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := pipeline.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			handleContextError(&pipeline.options, 5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := pipeline.decoder6(r, val1, val2, val3, val4, val5)
		if err != nil {
			handleContextError(&pipeline.options, 6, err).Respond(w, r)
			return
		}

		// Decode input (as seventh context)
		input, err := pipeline.decoder7(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5, val6, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipelineWithInput7 creates a handler with seven contexts and input as a pipeline stage
func HandlePipelineWithInput7[C1, C2, C3, C4, C5, C6, C7, T any](
	p Pipeline7[C1, C2, C3, C4, C5, C6, C7],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput7(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.decoder2(r, val1)
		if err != nil {
			handleContextError(&pipeline.options, 2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.decoder3(r, val1, val2)
		if err != nil {
			handleContextError(&pipeline.options, 3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.decoder4(r, val1, val2, val3)
		if err != nil {
			handleContextError(&pipeline.options, 4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := pipeline.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			handleContextError(&pipeline.options, 5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := pipeline.decoder6(r, val1, val2, val3, val4, val5)
		if err != nil {
			handleContextError(&pipeline.options, 6, err).Respond(w, r)
			return
		}

		// Decode seventh context
		val7, err := pipeline.decoder7(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			handleContextError(&pipeline.options, 7, err).Respond(w, r)
			return
		}

		// Decode input (as eighth context)
		input, err := pipeline.decoder8(r, val1, val2, val3, val4, val5, val6, val7)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5, val6, val7, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

// HandlePipelineWithInput8 creates a handler with eight contexts and input as a pipeline stage
func HandlePipelineWithInput8[C1, C2, C3, C4, C5, C6, C7, C8, T any](
	p Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8],
	inputDecoder func(r *http.Request) (T, error),
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7, val8 C8, input T) Responder,
) http.HandlerFunc {
	pipeline := NewPipelineWithInput8(p, inputDecoder)
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := pipeline.decoder1(r)
		if err != nil {
			handleContextError(&pipeline.options, 1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.decoder2(r, val1)
		if err != nil {
			handleContextError(&pipeline.options, 2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.decoder3(r, val1, val2)
		if err != nil {
			handleContextError(&pipeline.options, 3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.decoder4(r, val1, val2, val3)
		if err != nil {
			handleContextError(&pipeline.options, 4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := pipeline.decoder5(r, val1, val2, val3, val4)
		if err != nil {
			handleContextError(&pipeline.options, 5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := pipeline.decoder6(r, val1, val2, val3, val4, val5)
		if err != nil {
			handleContextError(&pipeline.options, 6, err).Respond(w, r)
			return
		}

		// Decode seventh context
		val7, err := pipeline.decoder7(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			handleContextError(&pipeline.options, 7, err).Respond(w, r)
			return
		}

		// Decode eighth context
		val8, err := pipeline.decoder8(r, val1, val2, val3, val4, val5, val6, val7)
		if err != nil {
			handleContextError(&pipeline.options, 8, err).Respond(w, r)
			return
		}

		// Decode input (as ninth context)
		input, err := pipeline.decoder9(r, val1, val2, val3, val4, val5, val6, val7, val8)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			handleInputError(&pipeline.options, err).Respond(w, r)
			return
		}

		// Call handler with request context
		res := handler(r.Context(), val1, val2, val3, val4, val5, val6, val7, val8, input)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		res.Respond(w, r)
	}
}

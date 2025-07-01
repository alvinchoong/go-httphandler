package pipeline

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alvinchoong/go-httphandler"
)

// ========== Pipeline options ==========

// PipelineOptions holds configurable options for pipelines
type PipelineOptions struct {
	// DecodeErrorHandler handles errors from context decoders
	DecodeErrorHandler func(stage int, err error) httphandler.Responder

	// InputErrorHandler handles errors from input decoders
	InputErrorHandler func(err error) httphandler.Responder
}

// ========== Flattened Pipeline Structures ==========

// Pipeline1 is a pipeline with one decoder type
type Pipeline1[C any] struct {
	decoder1 func(r *http.Request) (C, error)
	options  PipelineOptions
}

// Pipeline2 is a pipeline with two decoder types
type Pipeline2[C1, C2 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	options  PipelineOptions
}

// Pipeline3 is a pipeline with three decoder types
type Pipeline3[C1, C2, C3 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	options  PipelineOptions
}

// Pipeline4 is a pipeline with four decoder types
type Pipeline4[C1, C2, C3, C4 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	options  PipelineOptions
}

// Pipeline5 is a pipeline with five decoder types
type Pipeline5[C1, C2, C3, C4, C5 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error)
	options  PipelineOptions
}

// Pipeline6 is a pipeline with six decoder types
type Pipeline6[C1, C2, C3, C4, C5, C6 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error)
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error)
	options  PipelineOptions
}

// Pipeline7 is a pipeline with seven decoder types
type Pipeline7[C1, C2, C3, C4, C5, C6, C7 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error)
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error)
	decoder7 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error)
	options  PipelineOptions
}

// Pipeline8 is a pipeline with eight decoder types
type Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error)
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error)
	decoder7 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error)
	decoder8 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7) (C8, error)
	options  PipelineOptions
}

// ========== Factory methods for creating pipelines ==========

// New1 creates a pipeline with one decoder type
func New1[C any](
	decoder func(r *http.Request) (C, error),
	options ...func(*PipelineOptions),
) Pipeline1[C] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline1[C]{
		decoder1: decoder,
		options:  opts,
	}
}

// New2 creates a pipeline with two decoder types
func New2[C1, C2 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	options ...func(*PipelineOptions),
) Pipeline2[C1, C2] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline2[C1, C2]{
		decoder1: decoder1,
		decoder2: decoder2,
		options:  opts,
	}
}

// New3 creates a pipeline with three decoder types
func New3[C1, C2, C3 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	options ...func(*PipelineOptions),
) Pipeline3[C1, C2, C3] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline3[C1, C2, C3]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		options:  opts,
	}
}

// New4 creates a pipeline with four decoder types
func New4[C1, C2, C3, C4 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	options ...func(*PipelineOptions),
) Pipeline4[C1, C2, C3, C4] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline4[C1, C2, C3, C4]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		options:  opts,
	}
}

// New5 creates a pipeline with five decoder types
func New5[C1, C2, C3, C4, C5 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	options ...func(*PipelineOptions),
) Pipeline5[C1, C2, C3, C4, C5] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline5[C1, C2, C3, C4, C5]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		decoder5: decoder5,
		options:  opts,
	}
}

// New6 creates a pipeline with six decoder types
func New6[C1, C2, C3, C4, C5, C6 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error),
	options ...func(*PipelineOptions),
) Pipeline6[C1, C2, C3, C4, C5, C6] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline6[C1, C2, C3, C4, C5, C6]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		decoder5: decoder5,
		decoder6: decoder6,
		options:  opts,
	}
}

// New7 creates a pipeline with seven decoder types
func New7[C1, C2, C3, C4, C5, C6, C7 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error),
	decoder7 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error),
	options ...func(*PipelineOptions),
) Pipeline7[C1, C2, C3, C4, C5, C6, C7] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline7[C1, C2, C3, C4, C5, C6, C7]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		decoder5: decoder5,
		decoder6: decoder6,
		decoder7: decoder7,
		options:  opts,
	}
}

// New8 creates a pipeline with eight decoder types
func New8[C1, C2, C3, C4, C5, C6, C7, C8 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error),
	decoder7 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error),
	decoder8 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7) (C8, error),
	options ...func(*PipelineOptions),
) Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8] {
	opts := defaultOptions()
	for _, option := range options {
		option(&opts)
	}
	return Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		decoder5: decoder5,
		decoder6: decoder6,
		decoder7: decoder7,
		decoder8: decoder8,
		options:  opts,
	}
}

// ========== Error handling helper ==========

// errorResponder creates a Responder for decoder errors
func errorResponder(err error) httphandler.Responder {
	// Default to a 400 Bad Request for decode errors
	return &errorResponse{
		statusCode: http.StatusBadRequest,
		err:        err,
	}
}

// ========== Pipeline options helpers ==========

// defaultOptions returns the default PipelineOptions
func defaultOptions() PipelineOptions {
	return PipelineOptions{
		DecodeErrorHandler: func(stage int, err error) httphandler.Responder {
			return errorResponder(fmt.Errorf("context%d decode error: %w", stage, err))
		},
		InputErrorHandler: func(err error) httphandler.Responder {
			return errorResponder(fmt.Errorf("input decode error: %w", err))
		},
	}
}

// WithDecodeErrorHandler returns an option that sets a custom context error handler
func WithDecodeErrorHandler(handler func(stage int, err error) httphandler.Responder) func(*PipelineOptions) {
	return func(opts *PipelineOptions) {
		opts.DecodeErrorHandler = handler
	}
}

// WithInputErrorHandler returns an option that sets a custom input error handler
func WithInputErrorHandler(handler func(err error) httphandler.Responder) func(*PipelineOptions) {
	return func(opts *PipelineOptions) {
		opts.InputErrorHandler = handler
	}
}

// errorResponse implements the Responder interface for errors
type errorResponse struct {
	statusCode int
	err        error
}

// Respond implements the Responder interface
func (e *errorResponse) Respond(w http.ResponseWriter, r *http.Request) {
	http.Error(w, e.err.Error(), e.statusCode)
}

// ========== Pipeline builders ==========



// ========== Handler creator imports ==========

// See pipeline_handlers.go for handler implementations

// ========== Pipeline Handle methods ==========

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline1[C]) Handle(
	handler func(ctx context.Context, val C) httphandler.Responder,
) http.HandlerFunc {
	return Handle1(p, handler)
}

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline2[C1, C2]) Handle(
	handler func(ctx context.Context, val1 C1, val2 C2) httphandler.Responder,
) http.HandlerFunc {
	return Handle2(p, handler)
}

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline3[C1, C2, C3]) Handle(
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3) httphandler.Responder,
) http.HandlerFunc {
	return Handle3(p, handler)
}

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline4[C1, C2, C3, C4]) Handle(
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4) httphandler.Responder,
) http.HandlerFunc {
	return Handle4(p, handler)
}

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline5[C1, C2, C3, C4, C5]) Handle(
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5) httphandler.Responder,
) http.HandlerFunc {
	return Handle5(p, handler)
}

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline6[C1, C2, C3, C4, C5, C6]) Handle(
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6) httphandler.Responder,
) http.HandlerFunc {
	return Handle6(p, handler)
}

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline7[C1, C2, C3, C4, C5, C6, C7]) Handle(
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7) httphandler.Responder,
) http.HandlerFunc {
	return Handle7(p, handler)
}

// Handle creates a handler using this pipeline with no additional input decoder
func (p Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8]) Handle(
	handler func(ctx context.Context, val1 C1, val2 C2, val3 C3, val4 C4, val5 C5, val6 C6, val7 C7, val8 C8) httphandler.Responder,
) http.HandlerFunc {
	return Handle8(p, handler)
}

package httphandler

import (
	"fmt"
	"net/http"
)

// ========== Pipeline options ==========

// PipelineOptions holds configurable options for pipelines
type PipelineOptions struct {
	// DecodeErrorHandler handles errors from context decoders
	DecodeErrorHandler func(stage int, err error) Responder

	// InputErrorHandler handles errors from input decoders
	InputErrorHandler func(err error) Responder
}

// ========== Flattened Pipeline Structures ==========

// Pipeline1 is a pipeline with one context type
type Pipeline1[C any] struct {
	decoder1 func(r *http.Request) (C, error)
	options  PipelineOptions
}

// Pipeline2 is a pipeline with two context types
type Pipeline2[C1, C2 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	options  PipelineOptions
}

// Pipeline3 is a pipeline with three context types
type Pipeline3[C1, C2, C3 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	options  PipelineOptions
}

// Pipeline4 is a pipeline with four context types
type Pipeline4[C1, C2, C3, C4 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	options  PipelineOptions
}

// Pipeline5 is a pipeline with five context types
type Pipeline5[C1, C2, C3, C4, C5 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error)
	options  PipelineOptions
}

// Pipeline6 is a pipeline with six context types
type Pipeline6[C1, C2, C3, C4, C5, C6 any] struct {
	decoder1 func(r *http.Request) (C1, error)
	decoder2 func(r *http.Request, c1 C1) (C2, error)
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error)
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error)
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error)
	options  PipelineOptions
}

// Pipeline7 is a pipeline with seven context types
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

// Pipeline8 is a pipeline with eight context types
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

// NewPipeline1 creates a pipeline with one context type
func NewPipeline1[C any](
	decoder func(r *http.Request) (C, error),
	opts *PipelineOptions,
) Pipeline1[C] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
	}
	return Pipeline1[C]{
		decoder1: decoder,
		options:  options,
	}
}

// NewPipeline2 creates a pipeline with two context types
func NewPipeline2[C1, C2 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	opts *PipelineOptions,
) Pipeline2[C1, C2] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
	}
	return Pipeline2[C1, C2]{
		decoder1: decoder1,
		decoder2: decoder2,
		options:  options,
	}
}

// NewPipeline3 creates a pipeline with three context types
func NewPipeline3[C1, C2, C3 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	opts *PipelineOptions,
) Pipeline3[C1, C2, C3] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
	}
	return Pipeline3[C1, C2, C3]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		options:  options,
	}
}

// NewPipeline4 creates a pipeline with four context types
func NewPipeline4[C1, C2, C3, C4 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	opts *PipelineOptions,
) Pipeline4[C1, C2, C3, C4] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
	}
	return Pipeline4[C1, C2, C3, C4]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		options:  options,
	}
}

// NewPipeline5 creates a pipeline with five context types
func NewPipeline5[C1, C2, C3, C4, C5 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	opts *PipelineOptions,
) Pipeline5[C1, C2, C3, C4, C5] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
	}
	return Pipeline5[C1, C2, C3, C4, C5]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		decoder5: decoder5,
		options:  options,
	}
}

// NewPipeline6 creates a pipeline with six context types
func NewPipeline6[C1, C2, C3, C4, C5, C6 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error),
	opts *PipelineOptions,
) Pipeline6[C1, C2, C3, C4, C5, C6] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
	}
	return Pipeline6[C1, C2, C3, C4, C5, C6]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		decoder5: decoder5,
		decoder6: decoder6,
		options:  options,
	}
}

// NewPipeline7 creates a pipeline with seven context types
func NewPipeline7[C1, C2, C3, C4, C5, C6, C7 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error),
	decoder7 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error),
	opts *PipelineOptions,
) Pipeline7[C1, C2, C3, C4, C5, C6, C7] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
	}
	return Pipeline7[C1, C2, C3, C4, C5, C6, C7]{
		decoder1: decoder1,
		decoder2: decoder2,
		decoder3: decoder3,
		decoder4: decoder4,
		decoder5: decoder5,
		decoder6: decoder6,
		decoder7: decoder7,
		options:  options,
	}
}

// NewPipeline8 creates a pipeline with eight context types
func NewPipeline8[C1, C2, C3, C4, C5, C6, C7, C8 any](
	decoder1 func(r *http.Request) (C1, error),
	decoder2 func(r *http.Request, c1 C1) (C2, error),
	decoder3 func(r *http.Request, c1 C1, c2 C2) (C3, error),
	decoder4 func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	decoder5 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	decoder6 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error),
	decoder7 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error),
	decoder8 func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7) (C8, error),
	opts *PipelineOptions,
) Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8] {
	options := PipelineOptions{}
	if opts != nil {
		options = *opts
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
		options:  options,
	}
}

// ========== Error handling helper ==========

// errorResponder creates a Responder for decoder errors
func errorResponder(err error) Responder {
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
		DecodeErrorHandler: func(stage int, err error) Responder {
			return errorResponder(fmt.Errorf("context%d decode error: %w", stage, err))
		},
		InputErrorHandler: func(err error) Responder {
			return errorResponder(fmt.Errorf("input decode error: %w", err))
		},
	}
}

// WithDecodeErrorHandler returns an option that sets a custom context error handler
func WithDecodeErrorHandler(handler func(stage int, err error) Responder) func(*PipelineOptions) {
	return func(opts *PipelineOptions) {
		opts.DecodeErrorHandler = handler
	}
}

// WithInputErrorHandler returns an option that sets a custom input error handler
func WithInputErrorHandler(handler func(err error) Responder) func(*PipelineOptions) {
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

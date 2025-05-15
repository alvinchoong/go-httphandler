package httphandler

import (
	"fmt"
	"net/http"
)

// ========== Pipeline options ==========

// PipelineOptions holds configurable options for pipelines
type PipelineOptions struct {
	// ContextErrorHandler handles errors from context decoders
	ContextErrorHandler func(stage int, err error) Responder

	// InputErrorHandler handles errors from input decoders
	InputErrorHandler func(err error) Responder
}

// ========== Pipeline with one context type ==========

// Pipeline1 is a pipeline with one context type
type Pipeline1[C any] struct {
	decoder func(r *http.Request) (C, error)
	options PipelineOptions
}

// ========== Pipeline with two context types ==========

// Pipeline2 is a pipeline with two context types
type Pipeline2[C1, C2 any] struct {
	p1      Pipeline1[C1]
	decoder func(r *http.Request, c1 C1) (C2, error)
	options PipelineOptions
}

// ========== Pipeline with three context types ==========

// Pipeline3 is a pipeline with three context types
type Pipeline3[C1, C2, C3 any] struct {
	p2      Pipeline2[C1, C2]
	decoder func(r *http.Request, c1 C1, c2 C2) (C3, error)
	options PipelineOptions
}

// ========== Pipeline with four context types ==========

// Pipeline4 is a pipeline with four context types
type Pipeline4[C1, C2, C3, C4 any] struct {
	p3      Pipeline3[C1, C2, C3]
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
	options PipelineOptions
}

// ========== Pipeline with five context types ==========

// Pipeline5 is a pipeline with five context types
type Pipeline5[C1, C2, C3, C4, C5 any] struct {
	p4      Pipeline4[C1, C2, C3, C4]
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error)
	options PipelineOptions
}

// ========== Pipeline with six context types ==========

// Pipeline6 is a pipeline with six context types
type Pipeline6[C1, C2, C3, C4, C5, C6 any] struct {
	p5      Pipeline5[C1, C2, C3, C4, C5]
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error)
	options PipelineOptions
}

// ========== Pipeline with seven context types ==========

// Pipeline7 is a pipeline with seven context types
type Pipeline7[C1, C2, C3, C4, C5, C6, C7 any] struct {
	p6      Pipeline6[C1, C2, C3, C4, C5, C6]
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error)
	options PipelineOptions
}

// ========== Pipeline with eight context types ==========

// Pipeline8 is a pipeline with eight context types
type Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8 any] struct {
	p7      Pipeline7[C1, C2, C3, C4, C5, C6, C7]
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7) (C8, error)
	options PipelineOptions
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
		ContextErrorHandler: func(stage int, err error) Responder {
			return errorResponder(fmt.Errorf("context%d decode error: %w", stage, err))
		},
		InputErrorHandler: func(err error) Responder {
			return errorResponder(fmt.Errorf("input decode error: %w", err))
		},
	}
}

// WithContextErrorHandler returns an option that sets a custom context error handler
func WithContextErrorHandler(handler func(stage int, err error) Responder) func(*PipelineOptions) {
	return func(opts *PipelineOptions) {
		opts.ContextErrorHandler = handler
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

// NewPipeline1 creates a pipeline with one context type
func NewPipeline1[C any](
	decoder func(r *http.Request) (C, error),
	options ...func(*PipelineOptions),
) Pipeline1[C] {
	// Default options
	opts := defaultOptions()

	// Apply provided options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline1[C]{
		decoder: decoder,
		options: opts,
	}
}

// NewPipeline2 creates a pipeline with two context types
func NewPipeline2[C1, C2 any](
	p Pipeline1[C1],
	decoder func(r *http.Request, c1 C1) (C2, error),
	options ...func(*PipelineOptions),
) Pipeline2[C1, C2] {
	// Start with options from previous pipeline
	opts := p.options

	// Apply any new options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline2[C1, C2]{
		p1:      p,
		decoder: decoder,
		options: opts,
	}
}

// NewPipeline3 creates a pipeline with three context types
func NewPipeline3[C1, C2, C3 any](
	p Pipeline2[C1, C2],
	decoder func(r *http.Request, c1 C1, c2 C2) (C3, error),
	options ...func(*PipelineOptions),
) Pipeline3[C1, C2, C3] {
	// Start with options from previous pipeline
	opts := p.options

	// Apply any new options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline3[C1, C2, C3]{
		p2:      p,
		decoder: decoder,
		options: opts,
	}
}

// NewPipeline4 creates a pipeline with four context types
func NewPipeline4[C1, C2, C3, C4 any](
	p Pipeline3[C1, C2, C3],
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error),
	options ...func(*PipelineOptions),
) Pipeline4[C1, C2, C3, C4] {
	// Start with options from previous pipeline
	opts := p.options

	// Apply any new options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline4[C1, C2, C3, C4]{
		p3:      p,
		decoder: decoder,
		options: opts,
	}
}

// NewPipeline5 creates a pipeline with five context types
func NewPipeline5[C1, C2, C3, C4, C5 any](
	p Pipeline4[C1, C2, C3, C4],
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (C5, error),
	options ...func(*PipelineOptions),
) Pipeline5[C1, C2, C3, C4, C5] {
	// Start with options from previous pipeline
	opts := p.options

	// Apply any new options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline5[C1, C2, C3, C4, C5]{
		p4:      p,
		decoder: decoder,
		options: opts,
	}
}

// NewPipeline6 creates a pipeline with six context types
func NewPipeline6[C1, C2, C3, C4, C5, C6 any](
	p Pipeline5[C1, C2, C3, C4, C5],
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (C6, error),
	options ...func(*PipelineOptions),
) Pipeline6[C1, C2, C3, C4, C5, C6] {
	// Start with options from previous pipeline
	opts := p.options

	// Apply any new options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline6[C1, C2, C3, C4, C5, C6]{
		p5:      p,
		decoder: decoder,
		options: opts,
	}
}

// NewPipeline7 creates a pipeline with seven context types
func NewPipeline7[C1, C2, C3, C4, C5, C6, C7 any](
	p Pipeline6[C1, C2, C3, C4, C5, C6],
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (C7, error),
	options ...func(*PipelineOptions),
) Pipeline7[C1, C2, C3, C4, C5, C6, C7] {
	// Start with options from previous pipeline
	opts := p.options

	// Apply any new options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline7[C1, C2, C3, C4, C5, C6, C7]{
		p6:      p,
		decoder: decoder,
		options: opts,
	}
}

// NewPipeline8 creates a pipeline with eight context types
func NewPipeline8[C1, C2, C3, C4, C5, C6, C7, C8 any](
	p Pipeline7[C1, C2, C3, C4, C5, C6, C7],
	decoder func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7) (C8, error),
	options ...func(*PipelineOptions),
) Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8] {
	// Start with options from previous pipeline
	opts := p.options

	// Apply any new options
	for _, option := range options {
		option(&opts)
	}

	return Pipeline8[C1, C2, C3, C4, C5, C6, C7, C8]{
		p7:      p,
		decoder: decoder,
		options: opts,
	}
}

// ========== Handler creator imports ==========

// See pipeline_handlers.go for handler implementations

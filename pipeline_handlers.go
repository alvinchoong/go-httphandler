package httphandler

import (
	"context"
	"net/http"
)

// ========== Input as pipeline stage functions ==========

// NewPipelineWithInput1 creates a pipeline with one context type and input
func NewPipelineWithInput1[C, T any](
	p Pipeline1[C],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) Pipeline2[C, T] {
	return NewPipeline2(p, func(r *http.Request, c C) (T, error) {
		return inputDecoder(r)
	}, options...)
}

// NewPipelineWithInput2 creates a pipeline with two context types and input
func NewPipelineWithInput2[C1, C2, T any](
	p Pipeline2[C1, C2],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) Pipeline3[C1, C2, T] {
	return NewPipeline3(p, func(r *http.Request, c1 C1, c2 C2) (T, error) {
		return inputDecoder(r)
	}, options...)
}

// NewPipelineWithInput3 creates a pipeline with three context types and input
func NewPipelineWithInput3[C1, C2, C3, T any](
	p Pipeline3[C1, C2, C3],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) Pipeline4[C1, C2, C3, T] {
	return NewPipeline4(p, func(r *http.Request, c1 C1, c2 C2, c3 C3) (T, error) {
		return inputDecoder(r)
	}, options...)
}

// NewPipelineWithInput4 creates a pipeline with four context types and input
func NewPipelineWithInput4[C1, C2, C3, C4, T any](
	p Pipeline4[C1, C2, C3, C4],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) Pipeline5[C1, C2, C3, C4, T] {
	return NewPipeline5(p, func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4) (T, error) {
		return inputDecoder(r)
	}, options...)
}

// NewPipelineWithInput5 creates a pipeline with five context types and input
func NewPipelineWithInput5[C1, C2, C3, C4, C5, T any](
	p Pipeline5[C1, C2, C3, C4, C5],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) Pipeline6[C1, C2, C3, C4, C5, T] {
	return NewPipeline6(p, func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5) (T, error) {
		return inputDecoder(r)
	}, options...)
}

// NewPipelineWithInput6 creates a pipeline with six context types and input
func NewPipelineWithInput6[C1, C2, C3, C4, C5, C6, T any](
	p Pipeline6[C1, C2, C3, C4, C5, C6],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) Pipeline7[C1, C2, C3, C4, C5, C6, T] {
	return NewPipeline7(p, func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6) (T, error) {
		return inputDecoder(r)
	}, options...)
}

// NewPipelineWithInput7 creates a pipeline with seven context types and input
func NewPipelineWithInput7[C1, C2, C3, C4, C5, C6, C7, T any](
	p Pipeline7[C1, C2, C3, C4, C5, C6, C7],
	inputDecoder func(r *http.Request) (T, error),
	options ...func(*PipelineOptions),
) Pipeline8[C1, C2, C3, C4, C5, C6, C7, T] {
	return NewPipeline8(p, func(r *http.Request, c1 C1, c2 C2, c3 C3, c4 C4, c5 C5, c6 C6, c7 C7) (T, error) {
		return inputDecoder(r)
	}, options...)
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
		val1, err := pipeline.p1.decoder(r)
		if err != nil {
			pipeline.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode input (as second context)
		input, err := pipeline.decoder(r, val1)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			pipeline.options.InputErrorHandler(err).Respond(w, r)
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
		val1, err := pipeline.p2.p1.decoder(r)
		if err != nil {
			pipeline.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.p2.decoder(r, val1)
		if err != nil {
			pipeline.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode input (as third context)
		input, err := pipeline.decoder(r, val1, val2)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			pipeline.options.InputErrorHandler(err).Respond(w, r)
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
		val1, err := pipeline.p3.p2.p1.decoder(r)
		if err != nil {
			pipeline.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.p3.p2.decoder(r, val1)
		if err != nil {
			pipeline.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.p3.decoder(r, val1, val2)
		if err != nil {
			pipeline.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode input (as fourth context)
		input, err := pipeline.decoder(r, val1, val2, val3)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			pipeline.options.InputErrorHandler(err).Respond(w, r)
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
		val1, err := pipeline.p4.p3.p2.p1.decoder(r)
		if err != nil {
			pipeline.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.p4.p3.p2.decoder(r, val1)
		if err != nil {
			pipeline.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.p4.p3.decoder(r, val1, val2)
		if err != nil {
			pipeline.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.p4.decoder(r, val1, val2, val3)
		if err != nil {
			pipeline.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode input (as fifth context)
		input, err := pipeline.decoder(r, val1, val2, val3, val4)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			pipeline.options.InputErrorHandler(err).Respond(w, r)
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
		val1, err := pipeline.p5.p4.p3.p2.p1.decoder(r)
		if err != nil {
			pipeline.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.p5.p4.p3.p2.decoder(r, val1)
		if err != nil {
			pipeline.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.p5.p4.p3.decoder(r, val1, val2)
		if err != nil {
			pipeline.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.p5.p4.decoder(r, val1, val2, val3)
		if err != nil {
			pipeline.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := pipeline.p5.decoder(r, val1, val2, val3, val4)
		if err != nil {
			pipeline.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Decode input (as sixth context)
		input, err := pipeline.decoder(r, val1, val2, val3, val4, val5)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			pipeline.options.InputErrorHandler(err).Respond(w, r)
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
		val1, err := pipeline.p6.p5.p4.p3.p2.p1.decoder(r)
		if err != nil {
			pipeline.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.p6.p5.p4.p3.p2.decoder(r, val1)
		if err != nil {
			pipeline.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.p6.p5.p4.p3.decoder(r, val1, val2)
		if err != nil {
			pipeline.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.p6.p5.p4.decoder(r, val1, val2, val3)
		if err != nil {
			pipeline.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := pipeline.p6.p5.decoder(r, val1, val2, val3, val4)
		if err != nil {
			pipeline.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := pipeline.p6.decoder(r, val1, val2, val3, val4, val5)
		if err != nil {
			pipeline.options.ContextErrorHandler(6, err).Respond(w, r)
			return
		}

		// Decode input (as seventh context)
		input, err := pipeline.decoder(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			pipeline.options.InputErrorHandler(err).Respond(w, r)
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
		val1, err := pipeline.p7.p6.p5.p4.p3.p2.p1.decoder(r)
		if err != nil {
			pipeline.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := pipeline.p7.p6.p5.p4.p3.p2.decoder(r, val1)
		if err != nil {
			pipeline.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := pipeline.p7.p6.p5.p4.p3.decoder(r, val1, val2)
		if err != nil {
			pipeline.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := pipeline.p7.p6.p5.p4.decoder(r, val1, val2, val3)
		if err != nil {
			pipeline.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := pipeline.p7.p6.p5.decoder(r, val1, val2, val3, val4)
		if err != nil {
			pipeline.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := pipeline.p7.p6.decoder(r, val1, val2, val3, val4, val5)
		if err != nil {
			pipeline.options.ContextErrorHandler(6, err).Respond(w, r)
			return
		}

		// Decode seventh context
		val7, err := pipeline.p7.decoder(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			pipeline.options.ContextErrorHandler(7, err).Respond(w, r)
			return
		}

		// Decode input (as eighth context)
		input, err := pipeline.decoder(r, val1, val2, val3, val4, val5, val6, val7)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			pipeline.options.InputErrorHandler(err).Respond(w, r)
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
	// We can't create a NewPipelineWithInput8 since we only support up to Pipeline8
	// So we manually implement this function using the unified approach
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode first context
		val1, err := p.p7.p6.p5.p4.p3.p2.p1.decoder(r)
		if err != nil {
			p.options.ContextErrorHandler(1, err).Respond(w, r)
			return
		}

		// Decode second context
		val2, err := p.p7.p6.p5.p4.p3.p2.decoder(r, val1)
		if err != nil {
			p.options.ContextErrorHandler(2, err).Respond(w, r)
			return
		}

		// Decode third context
		val3, err := p.p7.p6.p5.p4.p3.decoder(r, val1, val2)
		if err != nil {
			p.options.ContextErrorHandler(3, err).Respond(w, r)
			return
		}

		// Decode fourth context
		val4, err := p.p7.p6.p5.p4.decoder(r, val1, val2, val3)
		if err != nil {
			p.options.ContextErrorHandler(4, err).Respond(w, r)
			return
		}

		// Decode fifth context
		val5, err := p.p7.p6.p5.decoder(r, val1, val2, val3, val4)
		if err != nil {
			p.options.ContextErrorHandler(5, err).Respond(w, r)
			return
		}

		// Decode sixth context
		val6, err := p.p7.p6.decoder(r, val1, val2, val3, val4, val5)
		if err != nil {
			p.options.ContextErrorHandler(6, err).Respond(w, r)
			return
		}

		// Decode seventh context
		val7, err := p.p7.decoder(r, val1, val2, val3, val4, val5, val6)
		if err != nil {
			p.options.ContextErrorHandler(7, err).Respond(w, r)
			return
		}

		// Decode eighth context
		val8, err := p.decoder(r, val1, val2, val3, val4, val5, val6, val7)
		if err != nil {
			p.options.ContextErrorHandler(8, err).Respond(w, r)
			return
		}

		// Decode input (which would be the ninth stage if we had Pipeline9)
		input, err := inputDecoder(r)
		if err != nil {
			// Use the InputErrorHandler for backward compatibility with tests
			p.options.InputErrorHandler(err).Respond(w, r)
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

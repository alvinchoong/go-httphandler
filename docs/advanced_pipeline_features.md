# Advanced Pipeline Features

This document outlines the design and implementation of advanced features for the request pipeline architecture. We have successfully implemented custom error handlers and extended the pipeline to support up to 8 contexts.

## 1. Custom Error Handlers

### Previous Implementation

Currently, the pipeline architecture uses an internal `errorResponder` function that converts errors to HTTP 400 Bad Request responses:

```go
func errorResponder(err error) Responder {
    // Default to a 400 Bad Request for decode errors
    return &errorResponse{
        statusCode: http.StatusBadRequest,
        err:        err,
    }
}
```

This approach is simple but doesn't allow for custom error handling strategies at different pipeline stages.

### Implemented Solution: Pipeline Options Pattern

We propose extending the pipeline types to include optional error handlers using a functional options pattern:

```go
// PipelineOptions holds configurable options for pipelines
type PipelineOptions struct {
    // DecodeErrorHandler handles errors from context decoders
    DecodeErrorHandler func(stage int, err error) Responder

    // InputErrorHandler handles errors from input decoders
    InputErrorHandler func(err error) Responder
}

// Pipeline1 with options
type Pipeline1[C any] struct {
    decoder1 func(r *http.Request) (C, error)
    options PipelineOptions
}
```

#### Builder Functions with Options

```go
// NewPipeline1 with functional options support
func NewPipeline1[C any](
    decoder func(r *http.Request) (C, error),
    options ...func(*PipelineOptions),
) Pipeline1[C] {
    opts := defaultOptions()
    for _, option := range options {
        option(&opts)
    }

    return Pipeline1[C]{
        decoder1: decoder,
        options: opts,
    }
}
```

#### Option Provider Functions

```go
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
```

#### Updated Handler Functions

```go
func HandlePipelineWithInput1[C, T any](
    p Pipeline1[C],
    inputDecoder func(r *http.Request) (T, error),
    handler func(ctx context.Context, val C, input T) Responder,
) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Decode context
        val, err := p.decoder1(r)
        if err != nil {
            p.options.DecodeErrorHandler(1, err).Respond(w, r)
            return
        }

        // Decode input
        input, err := inputDecoder(r)
        if err != nil {
            p.options.InputErrorHandler(err).Respond(w, r)
            return
        }

        // Call handler
        res := handler(r.Context(), val, input)
        if res == nil {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        res.Respond(w, r)
    }
}
```

#### Example Usage

```go
// Create a pipeline with custom error handling
tenantPipeline := httphandler.NewPipeline1(
    DecodeTenant,
    httphandler.WithDecodeErrorHandler(func(stage int, err error) httphandler.Responder {
        // Custom tenant-specific error handling
        if strings.Contains(err.Error(), "tenant not found") {
            return jsonresp.Error(nil, "Invalid tenant", http.StatusUnauthorized)
        }
        return jsonresp.Error(nil, err.Error(), http.StatusBadRequest)
    }),
)
```

### Benefits of This Approach

1. **Cleaner API**: No need to pass `nil` when no options are required
2. **Flexibility**: Allows for different error handling strategies at different pipeline stages
3. **Extensibility**: The functional options pattern can be extended for other customizations
4. **Type safety**: Maintains the type safety benefits of the current implementation
5. **Simplicity**: Users only need to provide options when they want to override defaults
6. **Composable**: Multiple options can be combined in a single function call

## Future Considerations

The current implementation provides a solid foundation for handling request processing with custom error handling. If additional advanced features are needed in the future, the PipelineOptions pattern can be extended to accommodate them while maintaining backward compatibility.

1. Define the `PipelineOptions` type with error handler fields
2. Update pipeline types to include options
3. Modify builder functions to accept options
4. Create option provider functions for error handlers
5. Update handler functions to use the custom error handlers
6. Write tests for custom error handling
7. Update documentation and examples

### Potential Future Enhancements

#### Logging and Metrics Integration

Instead of embedding logging and metrics directly into the pipeline options, consider implementation approaches that maintain separation of concerns:

1. **Middleware Pattern**: Allow users to wrap pipelines with their own middleware for logging/metrics
2. **Event Hooks**: Provide simple event hooks at key points in the pipeline
3. **Examples/Documentation**: Show how the library can be integrated with popular logging/metrics libraries

#### Conditional Pipeline Branching

If conditional branching is required in the future, consider these implementation approaches:

1. **Branch Handler Pattern**: Conditionally execute different handlers based on context values
2. **Branch Pipeline Pattern**: Allow different pipeline paths based on conditions
3. **Context-Based Routing**: Use decoded context to determine routing

## Backward Compatibility

All enhancements will maintain backward compatibility with existing code. Default implementations will behave exactly as the current code does, and users can opt into advanced features as needed.

## Implementation Status

### ✅ Completed

1. **Extended Pipeline Support**: Implemented pipeline architecture to support up to 8 contexts
2. **Custom Error Handlers**: Added support for stage-specific error handling

### 🚧 Future Considerations

1. **Logging and Metrics**: Integration with external logging and metrics systems
2. **Conditional Processing**: Support for branching pipelines based on context values
3. **Additional Configuration Options**: Further extension of the PipelineOptions pattern

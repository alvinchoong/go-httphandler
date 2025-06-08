# Request Pipeline Architecture Design

## Overview

This document outlines the design and implementation plan for adding a request pipeline architecture to the go-httphandler library. The pipeline architecture allows composing multiple request decoders into a processing chain that accumulates context as the request passes through each stage.

## Goals

1. Support middleware-like functionality while maintaining the library's simple, type-safe API
2. Enable composition of decoders for handling authentication, authorization, and data extraction
3. Allow handlers to receive accumulated context from previous pipeline stages
4. Maintain type safety throughout the pipeline
5. Keep the API intuitive and aligned with Go's idioms

## Core Concepts

### Pipeline Stages

Each pipeline stage is responsible for:
1. Extracting and validating a specific piece of information from the request
2. Accumulating that information into a context
3. Optionally responding early (e.g., for authentication failures)
4. Providing methods to add additional stages to the pipeline

### Decoders

Decoders are functions that extract and validate information from HTTP requests. They come in several variants:

1. **Initial Decoders**: Take only the request and return a specific type
   ```go
   // Generic initial decoder
   type ContextDecoder[C any] func(r *http.Request) (C, error)
   ```

2. **Contextual Decoders**: Take the request and context from previous stages
   ```go
   // Generic contextual decoder with one previous context
   type ContextualDecoder[C1, C2 any] func(r *http.Request, val1 C1) (C2, error)
   
   // Generic contextual decoder with two previous contexts
   type ContextualDecoder2[C1, C2, C3 any] func(r *http.Request, val1 C1, val2 C2) (C3, error)
   
   // And so on for more contexts
   ```

3. **Input Decoders**: Extract request-specific input data
   ```go
   // Generic input decoder
   type InputDecoder[T any] func(r *http.Request) (T, error)
   ```

### Pipeline Types

Due to Go's constraint that methods cannot have their own type parameters, we need to use a functional approach rather than an interface-based approach. We'll define concrete types for each pipeline depth and use free functions for pipeline operations:

```go
// Pipeline with one context type
type Pipeline1[C any] struct {
    decoder1 func(r *http.Request) (C, error)
    options  PipelineOptions
}

// Pipeline with two context types
type Pipeline2[C1, C2 any] struct {
    decoder1 func(r *http.Request) (C1, error)
    decoder2 func(r *http.Request, val1 C1) (C2, error)
    options  PipelineOptions
}

// Pipeline with three context types
type Pipeline3[C1, C2, C3 any] struct {
    decoder1 func(r *http.Request) (C1, error)
    decoder2 func(r *http.Request, val1 C1) (C2, error)
    decoder3 func(r *http.Request, val1 C1, val2 C2) (C3, error)
    options  PipelineOptions
}

// Extended to support Pipeline4, Pipeline5, Pipeline6, Pipeline7, and Pipeline8
```

## Design

### Pipeline Builders

```go
// Start a pipeline with one context type
func NewPipeline1[C any](
    decoder func(r *http.Request) (C, error),
    options ...func(*PipelineOptions),
) Pipeline1[C] {
    // Implementation
}

// Create a pipeline with two context types
func NewPipeline2[C1, C2 any](
    decoder1 func(r *http.Request) (C1, error),
    decoder2 func(r *http.Request, val1 C1) (C2, error),
    options ...func(*PipelineOptions),
) Pipeline2[C1, C2] {
    // Implementation
}

// Create a pipeline with three context types
func NewPipeline3[C1, C2, C3 any](
    decoder1 func(r *http.Request) (C1, error),
    decoder2 func(r *http.Request, val1 C1) (C2, error),
    decoder3 func(r *http.Request, val1 C1, val2 C2) (C3, error),
    options ...func(*PipelineOptions),
) Pipeline3[C1, C2, C3] {
    // Implementation
}

// Extended to support NewPipeline4 through NewPipeline8
```

### Handler Creation

```go
// Create a handler with one context and input
func HandlePipelineWithInput1[C, T any](
    p Pipeline1[C],
    inputDecoder func(r *http.Request) (T, error),
    handler func(ctx context.Context, val C, input T) Responder,
) http.HandlerFunc {
    // Implementation
}

// Create a handler with two contexts and input
func HandlePipelineWithInput2[C1, C2, T any](
    p Pipeline2[C1, C2],
    inputDecoder func(r *http.Request) (T, error),
    handler func(ctx context.Context, val1 C1, val2 C2, input T) Responder,
) http.HandlerFunc {
    // Implementation
}

// Extended to support HandleWithInput3 through HandleWithInput8
```

## Example Usage

This design allows developers to create pipelines with their own custom types:

```go
// Developer's custom types
type Tenant struct { /* ... */ }
type User struct { /* ... */ }
type Thing struct { /* ... */ }
type UpdateThingInput struct { /* ... */ }

// Developer's custom decoders
func DecodeTenant(r *http.Request) (Tenant, error) { /* ... */ }
func DecodeUser(r *http.Request, tenant Tenant) (User, error) { /* ... */ }
func DecodeThing(r *http.Request, tenant Tenant, user User) (Thing, error) { /* ... */ }

// Create pipeline stages
tenantPipeline := httphandler.NewPipeline1(DecodeTenant)
userPipeline := httphandler.NewPipeline2(DecodeTenant, DecodeUser)
thingPipeline := httphandler.NewPipeline3(DecodeTenant, DecodeUser, DecodeThing)

// Using the pipelines with different handler types
router.HandleFunc("/login", httphandler.HandleWithInput1(
    tenantPipeline,
    JSONBodyDecode[LoginInput],
    func(tenant Tenant, input LoginInput) httphandler.Responder {
        // Login handler with tenant context
        return jsonresp.Success(result)
    },
))

router.HandleFunc("/dashboard", httphandler.HandleWithInput2(
    userPipeline,
    QueryParamsDecode[DashboardParams],
    func(tenant Tenant, user User, params DashboardParams) httphandler.Responder {
        // Dashboard handler with tenant and user context
        return jsonresp.Success(result)
    },
))

router.HandleFunc("/api/things/{id}", httphandler.HandleWithInput3(
    thingPipeline,
    JSONBodyDecode[UpdateThingInput],
    func(tenant Tenant, user User, thing Thing, input UpdateThingInput) httphandler.Responder {
        // Update thing handler with all contexts
        return jsonresp.Success(result)
    },
))

## Implementation Status

### Core Implementation (Completed)

#### Implementation Details

The implementation of the request pipeline architecture uses Go's generics to provide type-safe request processing chains. Here's a detailed explanation of how the components work together:

##### 1. Pipeline Types

We've implemented concrete types for each pipeline depth (up to Pipeline8):

```go
// Pipeline1 is a pipeline with one context type
type Pipeline1[C any] struct {
    decoder func(r *http.Request) (C, error)
}

// Pipeline2 is a pipeline with two context types
type Pipeline2[C1, C2 any] struct {
    p1      Pipeline1[C1]
    decoder func(r *http.Request, c1 C1) (C2, error)
}

// Pipeline3 is a pipeline with three context types
type Pipeline3[C1, C2, C3 any] struct {
    p2      Pipeline2[C1, C2]
    decoder func(r *http.Request, c1 C1, c2 C2) (C3, error)
}

// Pipeline4 is a pipeline with four context types
type Pipeline4[C1, C2, C3, C4 any] struct {
    p3      Pipeline3[C1, C2, C3]
    decoder func(r *http.Request, c1 C1, c2 C2, c3 C3) (C4, error)
}
```

Each pipeline type stores its decoder function and a reference to the previous pipeline stage, allowing for proper chaining.

##### 2. Pipeline Builders

We've implemented free functions for creating and extending pipelines:

```go
// NewPipeline1 starts a pipeline with one context type
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

// NewPipeline2 creates a pipeline with two context types
func NewPipeline2[C1, C2 any](
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
        options: opts,
    }
}
```

And similar functions for WithContext3 and WithContext4. These functions make the API more intuitive and chainable.

##### 3. Handler Creation

We've implemented functions to create HTTP handlers from pipelines:

```go
// HandlePipelineWithInput1 creates a handler with one context and input
func HandlePipelineWithInput1[C, T any](
    p Pipeline1[C],
    inputDecoder func(r *http.Request) (T, error),
    handler func(ctx context.Context, val C, input T) Responder,
) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Decode context
        val, err := p.decoder1(r)
        if err != nil {
            errorResponder(fmt.Errorf("context decode error: %w", err)).Respond(w, r)
            return
        }

        // Decode input
        input, err := inputDecoder(r)
        if err != nil {
            errorResponder(fmt.Errorf("input decode error: %w", err)).Respond(w, r)
            return
        }

        // Call handler
        res := handler(ctx, input)
        if res == nil {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        res.Respond(w, r)
    }
}
```

And similar functions for HandleWithInput2, HandleWithInput3, and HandleWithInput4. These functions handle the orchestration of decoding contexts, decoding inputs, and calling the handler with the accumulated context.

##### 4. Error Handling

We've implemented a simple error handling mechanism for pipeline stages:

```go
// errorResponder creates a Responder for decoder errors
func errorResponder(err error) Responder {
    // Default to a 400 Bad Request for decode errors
    return &errorResponse{
        statusCode: http.StatusBadRequest,
        err:        err,
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
```

This ensures that errors from decoders are properly converted to HTTP responses.

#### Key Design Decisions

1. **Functional vs Interface-Based Approach**: Due to Go's constraint that methods cannot have their own type parameters, we use a functional approach rather than an interface-based approach. This allows for better type inference and more intuitive API.

2. **Concrete Types vs Generic Interfaces**: We use concrete types (Pipeline1, Pipeline2, etc.) instead of a single generic interface to enable better type checking and IDE support.

3. **Error Handling Strategy**: Errors from decoders are wrapped with context information and converted to HTTP responses. This ensures that client applications receive meaningful error messages.

4. **Composition over Inheritance**: The pipeline architecture uses composition (storing previous pipeline stages) rather than inheritance, following Go's idiomatic approach to code organization.

### Phase 1.5: Compile-Time Validation

1. Create a `pipeline_test.go` file with type-checking tests
2. Implement stub types and functions to validate the generic pipeline design
3. Ensure all pipeline combinations compile correctly
4. Verify type inference works as expected with the fluent API
5. Test edge cases like deeply nested pipeline chains

Example test file structure:
```go
func TestPipelineCompilation(t *testing.T) {
    // This test doesn't actually run assertions - it just ensures
    // the code compiles correctly with our generic types
    
    // Define test types
    type TestContext1 struct{}
    type TestContext2 struct{}
    type TestContext3 struct{}
    type TestInput struct{}
    
    // Define stub decoders
    decoder1 := func(r *http.Request) (TestContext1, error) {
        return TestContext1{}, nil
    }
    
    decoder2 := func(r *http.Request, val1 TestContext1) (TestContext2, error) {
        return TestContext2{}, nil
    }
    
    decoder3 := func(r *http.Request, val1 TestContext1, val2 TestContext2) (TestContext3, error) {
        return TestContext3{}, nil
    }
    
    inputDecoder := func(r *http.Request) (TestInput, error) {
        return TestInput{}, nil
    }
    
    // Create pipelines with different depth
    p1 := NewPipeline1(decoder1)
    p2 := NewPipeline2(decoder1, decoder2)
    p3 := NewPipeline3(decoder1, decoder2, decoder3)
    
    // Create handlers with various context depths
    _ = HandlePipelineWithInput1(p1, inputDecoder, func(ctx context.Context, val1 TestContext1, input TestInput) Responder {
        return nil // Stub implementation
    })
    
    _ = HandlePipelineWithInput2(p2, inputDecoder, func(ctx context.Context, val1 TestContext1, val2 TestContext2, input TestInput) Responder {
        return nil // Stub implementation
    })
    
    _ = HandlePipelineWithInput3(p3, inputDecoder, func(ctx context.Context, val1 TestContext1, val2 TestContext2, val3 TestContext3, input TestInput) Responder {
        return nil // Stub implementation
    })
    
    // If this function compiles, the test passes
}
```

### Phase 2: Standard Decoders (Completed)

#### Implemented Decoders

We've implemented a comprehensive set of standard decoders to handle common HTTP request elements:

##### 1. Header Decoders

```go
// HeaderValue returns a decoder that extracts a specific header value
func HeaderValue(name string) func(r *http.Request) (string, error) {
    return func(r *http.Request) (string, error) {
        value := r.Header.Get(name)
        if value == "" {
            return "", fmt.Errorf("header %q not found", name)
        }
        return value, nil
    }
}

// OptionalHeaderValue returns a decoder that extracts a header value if present
func OptionalHeaderValue(name string) func(r *http.Request) (string, error) {
    return func(r *http.Request) (string, error) {
        return r.Header.Get(name), nil
    }
}

// BearerToken extracts a bearer token from the Authorization header
func BearerToken(r *http.Request) (string, error) {
    auth := r.Header.Get("Authorization")
    if auth == "" {
        return "", fmt.Errorf("missing Authorization header")
    }
    
    if !strings.HasPrefix(auth, "Bearer ") {
        return "", fmt.Errorf("invalid Authorization header format, expected 'Bearer TOKEN'")
    }
    
    return strings.TrimPrefix(auth, "Bearer "), nil
}

// BasicAuth extracts username and password from Basic Auth
func BasicAuth(r *http.Request) (struct {
    Username string
    Password string
}, error) {
    username, password, ok := r.BasicAuth()
    if !ok {
        return struct {
            Username string
            Password string
        }{}, fmt.Errorf("missing or invalid Basic Auth")
    }
    
    return struct {
        Username string
        Password string
    }{
        Username: username,
        Password: password,
    }, nil
}
```

##### 2. Query Parameter Decoders

```go
// QueryParam returns a decoder that extracts a query parameter value
func QueryParam(name string) func(r *http.Request) (string, error) {
    return func(r *http.Request) (string, error) {
        value := r.URL.Query().Get(name)
        if value == "" {
            return "", fmt.Errorf("query parameter %q not found", name)
        }
        return value, nil
    }
}

// IntQueryParam returns a decoder that extracts a query parameter as an integer
func IntQueryParam(name string) func(r *http.Request) (int, error) {
    return ParsedQueryParam(name, strconv.Atoi)
}

// FloatQueryParam returns a decoder that extracts a query parameter as a float64
func FloatQueryParam(name string) func(r *http.Request) (float64, error) {
    return ParsedQueryParam(name, func(s string) (float64, error) {
        return strconv.ParseFloat(s, 64)
    })
}

// BoolQueryParam returns a decoder that extracts a query parameter as a boolean
func BoolQueryParam(name string) func(r *http.Request) (bool, error) {
    return ParsedQueryParam(name, strconv.ParseBool)
}
```

##### 3. Path Parameter Decoders

```go
// PathParam returns a decoder that extracts a path parameter value
// This works with Go 1.22's r.PathValue() method
func PathParam(name string) func(r *http.Request) (string, error) {
    return func(r *http.Request) (string, error) {
        value := r.PathValue(name)
        if value == "" {
            return "", fmt.Errorf("path parameter %q not found", name)
        }
        return value, nil
    }
}

// IntPathParam returns a decoder that extracts a path parameter as an integer
func IntPathParam(name string) func(r *http.Request) (int, error) {
    return ParsedPathParam(name, strconv.Atoi)
}
```

##### 4. Form Data Decoders

```go
// FormValue returns a decoder that extracts a form value
func FormValue(name string) func(r *http.Request) (string, error) {
    return func(r *http.Request) (string, error) {
        // Parse form data if not already parsed
        if err := r.ParseForm(); err != nil {
            return "", fmt.Errorf("failed to parse form: %w", err)
        }
        
        value := r.FormValue(name)
        if value == "" {
            return "", fmt.Errorf("form value %q not found", name)
        }
        return value, nil
    }
}

// FormFile extracts a file from a multipart form
func FormFile(name string) func(r *http.Request) (*MultipartFile, error) {
    return func(r *http.Request) (*MultipartFile, error) {
        // Parse multipart form with default max memory
        if err := r.ParseMultipartForm(32 << 20); err != nil {
            return nil, fmt.Errorf("failed to parse multipart form: %w", err)
        }
        
        file, header, err := r.FormFile(name)
        if err != nil {
            return nil, fmt.Errorf("form file %q not found: %w", name, err)
        }
        
        return &MultipartFile{
            File:   file,
            Header: header,
        }, nil
    }
}
```

##### 5. Body Decoders

```go
// JSONBody decodes the request body as JSON into a value of type T
func JSONBody[T any]() func(r *http.Request) (T, error) {
    return func(r *http.Request) (T, error) {
        var value T
        if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
            return value, fmt.Errorf("failed to decode JSON: %w", err)
        }
        return value, nil
    }
}
```

##### 6. Composite Decoders

```go
// Combine2 combines two decoders into a single decoder that returns a struct with both results
func Combine2[T1, T2 any](
    decoder1 func(r *http.Request) (T1, error),
    decoder2 func(r *http.Request) (T2, error),
) func(r *http.Request) (struct {
    V1 T1
    V2 T2
}, error) {
    return func(r *http.Request) (struct {
        V1 T1
        V2 T2
    }, error) {
        v1, err := decoder1(r)
        if err != nil {
            return struct {
                V1 T1
                V2 T2
            }{}, fmt.Errorf("decoder1 error: %w", err)
        }
        
        v2, err := decoder2(r)
        if err != nil {
            return struct {
                V1 T1
                V2 T2
            }{}, fmt.Errorf("decoder2 error: %w", err)
        }
        
        return struct {
            V1 T1
            V2 T2
        }{
            V1: v1,
            V2: v2,
        }, nil
    }
}
```

#### Usage Patterns

These standard decoders can be combined with the pipeline architecture to create powerful request processing flows. For example:

```go
// Create pipeline stages
tenantPipeline := httphandler.NewPipeline1(DecodeTenant)
userPipeline := httphandler.NewPipeline2(DecodeTenant, DecodeUser)

// Route that requires tenant and user authentication, and extracts an ID from path parameters
router.HandleFunc("GET /items/{id}", httphandler.HandlePipelineWithInput2(
    userPipeline,
    httphandler.IntPathParam("id"),
    func(ctx context.Context, tenant Tenant, user User, itemID int) httphandler.Responder {
        // Handler logic with tenant, user, and item ID
        return jsonresp.Success(item)
    },
))
```

This design provides a flexible and type-safe way to handle HTTP requests with multiple context values and input parameters.

### Phase 3: Advanced Features

1. Implement pipeline branching (conditional pipelines)
2. Add support for custom error handling at each pipeline stage
3. Implement logging and metrics integration
4. Add support for asynchronous pipeline stages (if needed)

### Phase 4: Documentation and Examples

1. Update library documentation with pipeline architecture
2. Create comprehensive examples:
   - Authentication & authorization example
   - Multi-tenant application example
   - Form handling with validation

## Technical Considerations

### Error Handling

Each pipeline stage will handle its own errors and can convert them to appropriate HTTP responses. This ensures errors are handled correctly at each stage.

### Type Safety

The design preserves type safety throughout the pipeline by using Go's generics. This prevents runtime type errors and improves developer experience.

### Performance

The pipeline implementation should minimize allocation and avoid unnecessary computation. Consider caching intermediate results when appropriate.

### Compatibility

The pipeline architecture should be compatible with the existing library API, allowing gradual adoption.

## Conclusion

The implemented request pipeline architecture provides a powerful way to handle complex HTTP request processing while maintaining the library's philosophy of making HTTP handlers more idiomatic and less error-prone. By using Go's generics, the architecture delivers several key benefits:

1. **Type Safety**: The entire pipeline chain maintains type safety, preventing runtime errors and enabling better IDE support.

2. **Separation of Concerns**: Each decoder is focused on extracting and validating one specific piece of information, making the code more maintainable.

3. **Composition**: Pipeline stages can be combined in various ways, creating reusable processing chains for different routes.

4. **Declarative Style**: The API allows for a declarative style of defining request processing, making the code more readable and less error-prone.

5. **Flexible Error Handling**: Errors at any stage are properly propagated and converted to appropriate HTTP responses.

The implementation includes concrete pipeline types (up to 4 levels deep), a comprehensive set of standard decoders, and helper functions for creating handlers. A complete example demonstrates how to use the architecture in a multi-tenant application context, showcasing the power and flexibility of this approach.

This architecture significantly enhances the go-httphandler library, providing developers with a more sophisticated yet intuitive way to handle HTTP requests while maintaining Go's idiomatic approach to code organization and error handling.

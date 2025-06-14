# go-httphandler

[![GoDoc](https://pkg.go.dev/badge/github.com/alvinchoong/go-httphandler)](https://pkg.go.dev/github.com/alvinchoong/go-httphandler)
[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard)
[![License](https://img.shields.io/github/license/alvinchoong/go-httphandler)](LICENSE)

A zero-dependency HTTP response handler for Go that makes writing HTTP handlers idiomatic and less error-prone.

## Features

- ⚡ **Zero Dependencies**: Built entirely on Go's standard library
- 📄 **Built-in Response Types**: Support for JSON, plain text, file downloads, and redirects
- 🛠️ **Fluent API**: Chain methods to customize responses with headers, cookies, and status codes
- 🔄 **Flexible Request Parsing**: Built-in JSON parsing with support for custom decoders
- 🧩 **Easily Extendable**: Create custom response types and request decoders
- 📝 **Integrated Logging**: Optional logging support for all response types
- 🔗 **Request Pipeline Architecture**: Type-safe middleware-like functionality using Go's generics

## Why go-httphandler?

Traditional Go HTTP handlers interact directly with `http.ResponseWriter`, which can lead to several common pitfalls:

```go
// Traditional approach - common pitfalls

// Pitfall 1: Headers must be set before writing the response
router.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    user := getUser(r.PathValue("id"))
    if user == nil {
        json.NewEncoder(w).Encode(map[string]string{
            "error": "User not found",
        })
        w.WriteHeader(http.StatusNotFound) // Bug: Too late! Headers can't be set after writing response
        return
    }
    json.NewEncoder(w).Encode(user)
})

// Pitfall 2: Missing returns cause code to continue executing
router.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    user := getUser(r.PathValue("id"))
    if user == nil {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "User not found",
        })
        // Missing return! Code continues executing...
    }
    
    // This will still execute!
    json.NewEncoder(w).Encode(user)
})

// go-httphandler approach - prevents both issues by design
router.HandleFunc("GET /users/{id}", httphandler.Handle(func(r *http.Request) httphandler.Responder {
    user := getUser(r.PathValue("id"))
    if user == nil {
        return jsonresp.Error(nil, &ErrorResponse{
            Code:    "NOT_FOUND",
            Message: "User not found",
        }, http.StatusNotFound)
    }
    return jsonresp.Success(user)
}))
```

## Installation

```bash
go get github.com/alvinchoong/go-httphandler
```

## Usage Examples

### Response Types

#### JSON Response

```go
func getUserHandler(r *http.Request) httphandler.Responder {
    user, err := getUser(r.PathValue("id"))
    if err != nil {
        return jsonresp.InternalServerError(err)
    }
    if user == nil {
        return jsonresp.Error(nil, &ErrorResponse{
            Code:    "NOT_FOUND",
            Message: "User not found",
        }, http.StatusNotFound)
    }
    return jsonresp.Success(user)
}

router.HandleFunc("GET /users/{id}", httphandler.Handle(getUserHandler))
```

#### File Response

```go
func downloadReportHandler(r *http.Request) httphandler.Responder {
    file := getReport()
    return downloadresp.Attachment(file, "report.pdf").
        WithContentType("application/pdf")
}
```

#### Redirect Response

```go
func redirectHandler(r *http.Request) httphandler.Responder {
    return httphandler.Redirect("/new-location", http.StatusTemporaryRedirect).
        WithCookie(&http.Cookie{Name: "session", Value: "123"})
}
```

#### Plain Text Response

```go
func healthCheckHandler(r *http.Request) httphandler.Responder {
    return plainresp.Success("OK").
        WithHeader("Cache-Control", "no-cache")
}
```

#### Response Customization

All responders support method chaining for customization:

```go
return jsonresp.Success(data).
    WithStatus(http.StatusAccepted).
    WithHeader("X-Custom-Header", "value").
    WithCookie(&http.Cookie{Name: "session", Value: "123"}).
    WithLogger(logger)
```

### Request Handling

#### JSON Request Parsing

```go
func createUserHandler(r *http.Request, input CreateUserInput) httphandler.Responder {
    if err := input.Validate(); err != nil {
        return jsonresp.Error(err, &ValidationError{
            Code:    "VALIDATION_ERROR",
            Message: "Invalid input",
            Details: err.Error(),
        }, http.StatusBadRequest)
    }
    
    user, err := createUser(input)
    if err != nil {
        return jsonresp.InternalServerError(err)
    }

    return jsonresp.Success(user)
}

router.HandleFunc("POST /users", httphandler.HandleWithInput(createUserHandler))
```

### Additional Examples

For more examples see:

- [Basic REST API Example](examples/basic/main.go) - Simple REST API implementation
- [Request Pipeline Architecture Example](examples/pipeline/main.go) - Demonstrates the pipeline architecture for multi-tenant applications

## Creating Custom Response Types

You can easily create your own response types by implementing the `Responder` interface.

### Custom CSV Responder

```go
// Define your custom responder
type CSVResponder struct {
    records    [][]string
    filename   string
    statusCode int
}

// Create a constructor
func NewCSVResponse(records [][]string, filename string) *CSVResponder {
    return &CSVResponder{
        records:    records,
        filename:   filename,
        statusCode: http.StatusOK,
    }
}

// Implement the Responder interface
func (res *CSVResponder) Respond(w http.ResponseWriter, r *http.Request) {
    // Set headers for CSV download
    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, res.filename))
    
    // Write status code
    w.WriteHeader(res.statusCode)
    
    // Write CSV
    writer := csv.NewWriter(w)
    if err := writer.WriteAll(res.records); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

// Usage example
func csvReportHandler(r *http.Request) httphandler.Responder {
    records := [][]string{
        {"Name", "Email", "Age"},
        {"John Doe", "john@example.com", "30"},
        {"Jane Doe", "jane@example.com", "28"},
    }

    return NewCSVResponse(records, "users.csv")
}
```

## Request Pipeline Architecture

The Request Pipeline Architecture provides a powerful way to chain HTTP request processing stages while maintaining type safety using Go's generics. This enables middleware-like functionality where each stage extracts and validates part of the request, accumulating context for your handler.

### Pipeline Types

```go
// Create a pipeline with one decoder type (e.g., Tenant)
tenantPipeline := httphandler.NewPipeline1(DecodeTenant)

// Add a second decoder type (e.g., User)
userPipeline := httphandler.NewPipeline2(tenantPipeline, DecodeUser)

// Add a third decoder type (e.g., Product)
productPipeline := httphandler.NewPipeline3(userPipeline, DecodeProduct)
```

### Handler Registration

```go
// Handler with tenant and user context plus input data
router.HandleFunc("POST /products", httphandler.HandlePipelineWithInput2(
    userPipeline,
    DecodeProductInput,
    func(ctx context.Context, tenant Tenant, user User, input ProductInput) httphandler.Responder {
        // Handler receives request context, decoded tenant, user, and input data
        // Your business logic here
        return jsonresp.Success(result)
    },
))

// Handler with tenant, user, and product context plus input data
router.HandleFunc("PUT /products/{id}", httphandler.HandlePipelineWithInput3(
    productPipeline,
    DecodeProductInput,
    func(ctx context.Context, tenant Tenant, user User, product Product, input ProductInput) httphandler.Responder {
        // Handler receives request context, decoded tenant, user, product, and input data
        // Your business logic here
        return jsonresp.Success(result)
    },
))
```

### Benefits

- **Type-Safe**: All context and input values are fully typed
- **Composable**: Build and reuse pipeline stages across routes
- **Error Handling**: Each stage can return errors that are automatically converted to HTTP responses
- **Separation of Concerns**: Validate request context before handler execution

For a complete example of the Request Pipeline Architecture, see [examples/pipeline/main.go](examples/pipeline/main.go) and its [documentation](examples/pipeline/README.md).

## Benchmarks

Performance comparison between standard Go HTTP handlers and `go-httphandler` (benchmarked on Apple M3 Pro):

```plain
BenchmarkJSONResponse/Go/StandardHTTPHandler                1167090      1013 ns/op      6118 B/op      18 allocs/op
BenchmarkJSONResponse/HTTPHandler/JSONResponse              1000000      1077 ns/op      6239 B/op      21 allocs/op
BenchmarkJSONRequest/Go/StandardHTTPHandler/JSONDecode       976839      1272 ns/op      6275 B/op      22 allocs/op
BenchmarkJSONRequest/HTTPHandler/JSONBodyDecode              926830      1250 ns/op      6379 B/op      26 allocs/op
```

Results show that `go-httphandler` adds a minimal and neglible overhead (< 100 nanoseconds) while providing significant safety and maintainability benefits.

You can validate these results on your system by running:

```bash
go test -bench=. -benchmem
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

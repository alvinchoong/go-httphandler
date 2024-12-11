# go-httphandler

[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard)
[![License](https://img.shields.io/github/license/alvinchoong/go-httphandler)](LICENSE)

A zero-dependency HTTP response handler for Go that makes writing HTTP handlers more idiomatic and less error-prone.

## Features

- ⚡ **Zero Dependencies**: Built entirely on Go's standard library
- 📄 **Built-in Response Types**: Support for JSON, plain text, file downloads, and redirects
- 🛠️ **Fluent API**: Chain methods to customize responses with headers, cookies, and status codes
- 🔄 **Flexible Request Parsing**: Built-in JSON parsing with support for custom decoders
- 🧩 **Easily Extendable**: Create custom response types and request decoders
- 📝 **Integrated Logging**: Optional logging support for all response types

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
        return jsonresp.Error(nil, "User not found", http.StatusNotFound)
    }
    return jsonresp.Success(user)
}))
```

## Installation

```bash
go get github.com/alvinchoong/go-httphandler
```

## Usage Examples

### Basic JSON Handler

```go
func getUserHandler(r *http.Request) httphandler.Responder {
    user := getUser(r.PathValue("id"))
    if user == nil {
        return jsonresp.Error(nil, "User not found", http.StatusNotFound)
    }
    return jsonresp.Success(user)
}

router.HandleFunc("GET /users/{id}", httphandler.Handle(getUserHandler))
```

### Handler with Request Body Parsing

```go
func createUserHandler(r *http.Request, input CreateUserInput) httphandler.Responder {
    user, err := createUser(input)
    if err != nil {
        return jsonresp.Error(err, "Failed to create user", http.StatusBadRequest)
    }
    return jsonresp.Success(user)
}

router.HandleFunc("POST /users", httphandler.HandleWithInput(createUserHandler))
```

### File Download

```go
func downloadReportHandler(r *http.Request) httphandler.Responder {
    file, err := getReport()
    if err != nil {
        return jsonresp.Error(err, "Report not found", http.StatusNotFound)
    }
    
    return downloadresp.Attachment(file, "report.pdf").
        WithContentType("application/pdf").
        WithLogger(logger)
}
```

### Redirect Response

```go
func redirectHandler(r *http.Request) httphandler.Responder {
    return httphandler.Redirect("/new-location", http.StatusTemporaryRedirect).
        WithCookie(&http.Cookie{Name: "session", Value: "123"})
}
```

### Plain Text Response

```go
func healthCheckHandler(r *http.Request) httphandler.Responder {
    return plainresp.Success("OK").
        WithHeader("Cache-Control", "no-cache")
}
```

### Additional Examples

For more examples including a full REST API implementation see [examples/main.go](examples/main.go)

## Response Customization

All responders support method chaining for customization:

```go
return jsonresp.Success(data).
    WithStatus(http.StatusAccepted).
    WithHeader("X-Custom-Header", "value").
    WithCookie(&http.Cookie{Name: "session", Value: "123"}).
    WithLogger(logger)
```

## Error Handling

The library provides consistent error handling across all response types:

```go
// Standard error response
return jsonresp.Error(err, "Invalid input", http.StatusBadRequest)

// Internal server error with logging
return jsonresp.InternalServerError(err).WithLogger(logger)

// Plain text error
return plainresp.Error(err, "Not Found", http.StatusNotFound)
```

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
    return NewCSVResponse(records, "users.csv").
        WithHeader("Cache-Control", "no-cache")
}
```

## Benchmarks

Performance comparison between standard Go HTTP handlers and go-httphandler (benchmarked on Apple M3 Pro):

```plain
BenchmarkJSONResponse/Go/StandardHTTPHandler                1164493    1017 ns/op    6118 B/op    18 allocs/op
BenchmarkJSONResponse/HTTPHandler/JSONResponse             1000000    1091 ns/op    6294 B/op    21 allocs/op
BenchmarkJSONRequestResponse/Go/StandardHTTPHandlerWithInput       978330    1201 ns/op    6275 B/op    22 allocs/op
BenchmarkJSONRequestResponse/HTTPHandler/JSONRequestResponse       901645    1270 ns/op    6379 B/op    26 allocs/op
```

Results show that `go-httphandler` adds a minimal overhead (less than 75 nanoseconds) while providing significant safety and maintainability benefits. This overhead is negligible in real-world applications where network latency (typically milliseconds) and business logic are the primary performance factors.

You can run the benchmarks yourself:

```bash
go test -bench=. -benchmem
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

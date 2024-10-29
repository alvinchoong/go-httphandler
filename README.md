# go-httphandler

[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard)
[![License](https://img.shields.io/github/license/alvinchoong/go-httphandler)](LICENSE)

A zero-dependency HTTP handler wrapper for Go that enables safer, more idiomatic HTTP handlers by returning values instead of writing directly to `http.ResponseWriter`.

## Introduction

`go-httphandler` simplifies HTTP handler creation by letting you return values instead of writing directly to `http.ResponseWriter`. This design helps prevent common mistakes, such as inadvertently allowing code to continue executing after sending a response.

## Key Features

- 🚀 **Zero Dependencies:** Built purely on Go's standard library.
- 🔄 **Automatic JSON Unmarshalling:** Simplify input handling with automatic JSON parsing into structured types.
- 🛠 **Standardized Error Handling:** Consistent error responses across all handlers.
- 📄 **Multiple Response Types:** Support for JSON, plain text, file downloads, and redirects out-of-the-box.

## Why go-httphandler?

Traditional Go HTTP handlers interact directly with `http.ResponseWriter`, which can lead to several common pitfalls:

```go
// Pitfall 1: Headers must be set before writing the response
http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    user, exists := users[id]
    if !exists {
        json.NewEncoder(w).Encode(map[string]string{
            "error": "User not found",
        })
        w.WriteHeader(http.StatusNotFound) // Bug: Too late! Headers can't be set after writing response
        return
    }
    json.NewEncoder(w).Encode(user)
})

// Pitfall 2: Missing returns cause code to continue executing
http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    user, exists := users[id]
    if !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "User not found",
        })
        // Missing return! Code continues executing...
    }
    
    // This will still execute if user doesn't exist
    json.NewEncoder(w).Encode(user)
})

// go-httphandler approach - prevents both issues by design
router.HandleFunc("GET /users/{id}", httphandler.Handle(func(r *http.Request) httphandler.Responder {
    id := r.PathValue("id")
    user, exists := users[id]
    if !exists {
        return jsonresp.Error(nil, "User not found", http.StatusNotFound)
    }
    return jsonresp.Success(&user)
}))
```

## Installation

```sh
go get github.com/alvinchoong/go-httphandler
```

## Basic Usage

```go
// Basic handler
router.HandleFunc("GET /users/{id}", httphandler.Handle(func(r *http.Request) httphandler.Responder {
    user := getUser(r.PathValue("id"))
    return jsonresp.Success(&user)
}))

// Handler with input parsing
router.HandleFunc("POST /users", httphandler.HandleWithInput(func(r *http.Request, input UserInput) httphandler.Responder {
    user := createUser(input)
    return jsonresp.Success(&user).WithStatus(http.StatusCreated)
}))

// Error handling
return jsonresp.Error(err, "User not found", http.StatusBadRequest)

return jsonresp.InternalServerError(err)

// File download
return downloadresp.Attachment(fileReader, "report.pdf").
    WithContentType("application/pdf").
    WithLogger(logger)

// Redirect
return httphandler.Redirect("/new-path", http.StatusTemporaryRedirect)
```

For complete examples including a full REST API implementation, see [examples/main.go](examples/main.go)

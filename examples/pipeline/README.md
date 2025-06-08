# Request Pipeline Architecture Example

This example demonstrates the Request Pipeline Architecture in go-httphandler, which enables middleware-like functionality for HTTP handlers while maintaining type safety through Go's generics.

## Concepts Demonstrated

1. **Pipeline Chains**: Building a processing pipeline that accumulates context through multiple stages
2. **Context Decoders**: Extracting and validating data from HTTP requests (headers, path params, etc.)
3. **Multi-tenant Authentication**: Validating tenant and user context before processing requests
4. **Handler Composition**: Creating handlers that receive multiple context values

## Architecture Overview

This example implements a multi-tenant product management API with the following pipeline stages:

```
DecodeTenant → DecodeUser → DecodeProduct → Handle Request
```

Each stage adds context that is available to subsequent stages:

1. **Tenant Context**: Extracted from the `X-Tenant-ID` header
2. **User Context**: Extracted from the `Authorization` header, validated against the tenant
3. **Product Context**: Extracted from the URL path parameter, validated against the tenant

## Example Endpoints

| Method | Path             | Required Context        | Description                   |
|--------|------------------|-------------------------|-------------------------------|
| GET    | /products        | Tenant, User            | List products for tenant      |
| GET    | /products/{id}   | Tenant, User, Product   | Get specific product          |
| POST   | /products        | Tenant, User, Input     | Create new product            |
| PUT    | /products/{id}   | Tenant, User, Product, Input | Update existing product  |

## Testing the API

Run the server:

```bash
go run examples/pipeline/main.go
```

Use curl or any HTTP client with these headers:

```
X-Tenant-ID: t1            # or t2
Authorization: Bearer u1    # or u2 or u3
```

Example requests:

```bash
# List products for tenant t1
curl -H "X-Tenant-ID: t1" -H "Authorization: Bearer u1" http://localhost:8080/products

# Get product p1
curl -H "X-Tenant-ID: t1" -H "Authorization: Bearer u1" http://localhost:8080/products/p1

# Create new product (requires admin role)
curl -X POST -H "X-Tenant-ID: t1" -H "Authorization: Bearer u1" \
     -H "Content-Type: application/json" \
     -d '{"name":"New Product","price":299.99}' \
     http://localhost:8080/products

# Update product p1 (requires admin role)
curl -X PUT -H "X-Tenant-ID: t1" -H "Authorization: Bearer u1" \
     -H "Content-Type: application/json" \
     -d '{"name":"Updated Product","price":399.99}' \
     http://localhost:8080/products/p1
```

## Key Implementation Details

### Pipeline Definition

The pipeline is built using the `NewPipeline1`, `NewPipeline2`, and `NewPipeline3` functions:

```go
tenantPipeline := httphandler.NewPipeline1(DecodeTenant)
userPipeline := httphandler.NewPipeline2(tenantPipeline, DecodeUser)
productPipeline := httphandler.NewPipeline3(userPipeline, DecodeProduct)
```

### Handler Registration

Handlers are registered using the `HandlePipelineWithInput2` and `HandlePipelineWithInput3` functions, which provide accumulated context and input data:

```go
router.HandleFunc("GET /products/{id}", httphandler.HandlePipelineWithInput3(
    productPipeline,
    func(r *http.Request) (struct{}, error) { return struct{}{}, nil },
    func(ctx context.Context, tenant Tenant, user User, product Product, _ struct{}) httphandler.Responder {
        return GetProduct(tenant, user, product)
    },
))
```

### Error Handling

Each pipeline stage can return an error, which is automatically converted to an appropriate HTTP response. This ensures that validation failures result in consistent error responses without additional code.

## Benefits of the Pipeline Architecture

1. **Type Safety**: All context values and handler parameters are type-checked at compile time
2. **Separation of Concerns**: Each stage focuses on a specific responsibility
3. **Reusability**: Pipeline stages can be composed and reused across multiple routes
4. **Centralized Validation**: Input validation happens before your handler code executes
5. **Consistent Error Handling**: Pipeline errors are handled uniformly across all handlers

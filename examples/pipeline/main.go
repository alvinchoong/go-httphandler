package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/jsonresp"
)

// ========== Domain Types ==========

// Tenant represents a tenant in a multi-tenant application
type Tenant struct {
	ID   string
	Name string
}

// Response represents a standard API response
type Response struct {
	Tenant   string    `json:"tenant,omitempty"`
	User     string    `json:"user,omitempty"`
	Product  *Product  `json:"product,omitempty"`
	Products []Product `json:"products,omitempty"`
	Message  string    `json:"message,omitempty"`
}

// User represents a user within a tenant
type User struct {
	ID       string
	TenantID string
	Username string
	Role     string
}

// Product represents a product in the system
type Product struct {
	ID       string
	TenantID string
	Name     string
	Price    float64
}

// ProductInput represents the input for creating/updating a product
type ProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ========== In-memory data stores (for example purposes) ==========

var tenants = map[string]Tenant{
	"t1": {ID: "t1", Name: "Acme Corp"},
	"t2": {ID: "t2", Name: "Globex Inc"},
}

var users = map[string]User{
	"u1": {ID: "u1", TenantID: "t1", Username: "alice", Role: "admin"},
	"u2": {ID: "u2", TenantID: "t1", Username: "bob", Role: "user"},
	"u3": {ID: "u3", TenantID: "t2", Username: "carol", Role: "admin"},
}

var products = map[string]Product{
	"p1": {ID: "p1", TenantID: "t1", Name: "Product A", Price: 99.99},
	"p2": {ID: "p2", TenantID: "t1", Name: "Product B", Price: 149.99},
	"p3": {ID: "p3", TenantID: "t2", Name: "Product C", Price: 199.99},
}

// ========== Decoder functions ==========

// DecodeTenant extracts and validates tenant from the hostname
func DecodeTenant(r *http.Request) (Tenant, error) {
	// In a real app, we'd extract tenant ID from subdomain, path, or header
	// For this example, we'll use a custom header
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		return Tenant{}, fmt.Errorf("missing tenant ID")
	}

	tenant, found := tenants[tenantID]
	if !found {
		return Tenant{}, fmt.Errorf("tenant not found: %s", tenantID)
	}

	return tenant, nil
}

// DecodeUser extracts and validates user from auth header
func DecodeUser(r *http.Request, tenant Tenant) (User, error) {
	// In a real app, we'd decode a JWT token or check session
	// For this example, we'll use a simple auth header
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return User{}, fmt.Errorf("invalid authorization header")
	}

	userID := strings.TrimPrefix(auth, "Bearer ")
	user, found := users[userID]
	if !found {
		return User{}, fmt.Errorf("user not found: %s", userID)
	}

	// Ensure user belongs to the tenant
	if user.TenantID != tenant.ID {
		return User{}, fmt.Errorf("user does not belong to tenant")
	}

	return user, nil
}

// DecodeProduct extracts and validates product from path parameter
func DecodeProduct(r *http.Request, tenant Tenant, user User) (Product, error) {
	// Extract product ID from path
	productID := r.PathValue("id")
	if productID == "" {
		return Product{}, fmt.Errorf("missing product ID")
	}

	product, found := products[productID]
	if !found {
		return Product{}, fmt.Errorf("product not found: %s", productID)
	}

	// Ensure product belongs to the tenant
	if product.TenantID != tenant.ID {
		return Product{}, fmt.Errorf("product does not belong to tenant")
	}

	return product, nil
}

// DecodeProductInput parses product input from request body
func DecodeProductInput(r *http.Request) (ProductInput, error) {
	var input ProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return ProductInput{}, fmt.Errorf("invalid request body: %w", err)
	}

	// Validate input
	if input.Name == "" {
		return ProductInput{}, fmt.Errorf("name is required")
	}
	if input.Price <= 0 {
		return ProductInput{}, fmt.Errorf("price must be greater than zero")
	}

	return input, nil
}

// ========== Handler functions ==========

// ListProducts returns all products for a tenant
func ListProducts(tenant Tenant, user User) httphandler.Responder {
	// Get products for this tenant
	var tenantProducts []Product
	for _, product := range products {
		if product.TenantID == tenant.ID {
			tenantProducts = append(tenantProducts, product)
		}
	}

	return jsonresp.Success[Response](&Response{
		Tenant:   tenant.Name,
		User:     user.Username,
		Products: tenantProducts,
	})
}

// GetProduct returns a specific product
func GetProduct(tenant Tenant, user User, product Product) httphandler.Responder {
	return jsonresp.Success[Response](&Response{
		Tenant:  tenant.Name,
		User:    user.Username,
		Product: &product,
	})
}

// CreateProduct creates a new product
func CreateProduct(tenant Tenant, user User, input ProductInput) httphandler.Responder {
	// Check if user has admin role
	if user.Role != "admin" {
		return jsonresp.Error[string](fmt.Errorf("permission denied"), nil, http.StatusForbidden)
	}

	// Create new product
	productID := fmt.Sprintf("p%d", len(products)+1)
	product := Product{
		ID:       productID,
		TenantID: tenant.ID,
		Name:     input.Name,
		Price:    input.Price,
	}

	// Save product (in real app, this would be a database operation)
	products[productID] = product

	return jsonresp.Success[Response](&Response{
		Message: "product created",
		Product: &product,
	}).WithStatus(http.StatusCreated)
}

// UpdateProduct updates an existing product
func UpdateProduct(tenant Tenant, user User, product Product, input ProductInput) httphandler.Responder {
	// Check if user has admin role
	if user.Role != "admin" {
		return jsonresp.Error[string](fmt.Errorf("permission denied"), nil, http.StatusForbidden)
	}

	// Update product
	product.Name = input.Name
	product.Price = input.Price

	// Save updated product (in real app, this would be a database operation)
	products[product.ID] = product

	return jsonresp.Success[Response](&Response{
		Message: "product updated",
		Product: &product,
	})
}

func main() {
	// Create pipeline stages with the new flattened structure
	// No need to pass options when default error handling is sufficient
	userPipeline := httphandler.NewPipeline2(DecodeTenant, DecodeUser)
	productPipeline := httphandler.NewPipeline3(DecodeTenant, DecodeUser, DecodeProduct)

	// Set up router
	router := http.NewServeMux()

	// Route: List products (requires tenant and user)
	router.HandleFunc("GET /products", httphandler.HandlePipelineWithInput2(
		userPipeline,
		func(r *http.Request) (struct{}, error) { return struct{}{}, nil },
		func(ctx context.Context, tenant Tenant, user User, _ struct{}) httphandler.Responder {
			return ListProducts(tenant, user)
		},
	))

	// Route: Get product (requires tenant, user, and product)
	router.HandleFunc("GET /products/{id}", httphandler.HandlePipelineWithInput3(
		productPipeline,
		func(r *http.Request) (struct{}, error) { return struct{}{}, nil },
		func(ctx context.Context, tenant Tenant, user User, product Product, _ struct{}) httphandler.Responder {
			return GetProduct(tenant, user, product)
		},
	))

	// Route: Create product (requires tenant, user, and input)
	router.HandleFunc("POST /products", httphandler.HandlePipelineWithInput2(
		userPipeline,
		DecodeProductInput,
		func(ctx context.Context, tenant Tenant, user User, input ProductInput) httphandler.Responder {
			return CreateProduct(tenant, user, input)
		},
	))

	// Route: Update product (requires tenant, user, product, and input)
	router.HandleFunc("PUT /products/{id}", httphandler.HandlePipelineWithInput3(
		productPipeline,
		DecodeProductInput,
		func(ctx context.Context, tenant Tenant, user User, product Product, input ProductInput) httphandler.Responder {
			return UpdateProduct(tenant, user, product, input)
		},
	))

	// Start server
	port := ":8080"
	slog.Info(fmt.Sprintf("Pipeline example server starting on %s", port))
	slog.Info("To test the API, use these headers:")
	slog.Info("X-Tenant-ID: t1 or t2")
	slog.Info("Authorization: Bearer u1 or Bearer u2 or Bearer u3")
	if err := http.ListenAndServe(port, router); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}

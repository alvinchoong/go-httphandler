package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/alvinchoong/go-httphandler"
	"github.com/alvinchoong/go-httphandler/downloadresp"
	"github.com/alvinchoong/go-httphandler/jsonresp"
	"github.com/alvinchoong/go-httphandler/plainresp"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

type UserInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Store users in memory for demo purposes
var users = make(map[string]User)

func main() {
	router := http.NewServeMux()

	// Create User
	router.HandleFunc("POST /users", httphandler.HandleWithInput(createUser))

	// Get User by ID
	router.HandleFunc("GET /users/{id}", httphandler.Handle(getUser))

	// List Users
	router.HandleFunc("GET /users", httphandler.Handle(listUsers))

	// Update User by ID
	router.HandleFunc("PUT /users/{id}", httphandler.HandleWithInput(updateUser))

	// Delete User by ID
	router.HandleFunc("DELETE /users/{id}", httphandler.Handle(deleteUser))

	// Download users as CSV
	router.HandleFunc("GET /users/download", httphandler.Handle(downloadUsers))

	// Redirect example
	router.HandleFunc("GET /redirect", httphandler.Handle(redirectExample))

	// Plain text response example
	router.HandleFunc("GET /hello", httphandler.Handle(helloWorld))

	// Start server
	slog.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}

func createUser(r *http.Request, input UserInput) httphandler.Responder {
	if err := validate(input); err != nil {
		return jsonresp.Error(err, err.Error(), http.StatusBadRequest)
	}

	user := User{
		ID:        fmt.Sprintf("%d", time.Now().Unix()),
		Name:      input.Name,
		Age:       input.Age,
		CreatedAt: time.Now(),
	}
	users[user.ID] = user

	return jsonresp.Success(&user).WithStatus(http.StatusCreated)
}

func getUser(r *http.Request) httphandler.Responder {
	id := r.PathValue("id")
	user, exists := users[id]
	if !exists {
		return jsonresp.Error(nil, "User not found", http.StatusNotFound)
	}

	return jsonresp.Success(&user)
}

func listUsers(r *http.Request) httphandler.Responder {
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	return jsonresp.Success(&userList)
}

func updateUser(r *http.Request, input UserInput) httphandler.Responder {
	id := r.PathValue("id")
	user, exists := users[id]
	if !exists {
		return jsonresp.Error(nil, "User not found", http.StatusNotFound)
	}

	if err := validate(input); err != nil {
		return jsonresp.Error(err, err.Error(), http.StatusBadRequest)
	}

	user.Name = input.Name
	user.Age = input.Age
	users[id] = user

	return jsonresp.Success(&user)
}

func deleteUser(r *http.Request) httphandler.Responder {
	id := r.PathValue("id")
	if _, exists := users[id]; !exists {
		return jsonresp.Error(nil, "User not found", http.StatusNotFound)
	}

	delete(users, id)
	return jsonresp.Success[User](nil).WithStatus(http.StatusNoContent)
}

func downloadUsers(r *http.Request) httphandler.Responder {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()

		csvWriter := csv.NewWriter(pw)
		defer csvWriter.Flush()

		// Write header
		csvWriter.Write([]string{"ID", "Name", "Age", "Created At"})

		for _, user := range users {
			if err := csvWriter.Write([]string{
				user.ID,
				user.Name,
				strconv.Itoa(user.Age),
				user.CreatedAt.Format(time.RFC3339),
			}); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	return downloadresp.Attachment(pr, "users.csv").
		WithContentType("text/csv")
}

func redirectExample(r *http.Request) httphandler.Responder {
	return httphandler.Redirect("https://google.com", http.StatusSeeOther)
}

func helloWorld(r *http.Request) httphandler.Responder {
	return plainresp.Success("Hello, World!")
}

func validate(input UserInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Age <= 0 {
		return fmt.Errorf("age must be > 0")
	}
	return nil
}

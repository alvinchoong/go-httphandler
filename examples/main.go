package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"httphandler"
	"httphandler/downloadresp"
	"httphandler/jsonresp"
	"httphandler/plainresp"
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

var users = make(map[string]User)

func main() {
	router := http.NewServeMux()

	// Create User
	router.HandleFunc("POST /users", httphandler.HandleWithInput(func(r *http.Request, input UserInput) httphandler.Responder {
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
	}))

	// Get User by ID
	router.HandleFunc("GET /users/{id}", httphandler.Handle(func(r *http.Request) httphandler.Responder {
		id := r.PathValue("id")
		user, exists := users[id]
		if !exists {
			return jsonresp.Error(nil, "User not found", http.StatusNotFound)
		}

		return jsonresp.Success(&user)
	}))

	// List Users
	router.HandleFunc("GET /users", httphandler.Handle(func(r *http.Request) httphandler.Responder {
		userList := make([]User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}
		return jsonresp.Success(&userList)
	}))

	// Update User by ID
	router.HandleFunc("PUT /users/{id}", httphandler.HandleWithInput(func(r *http.Request, input UserInput) httphandler.Responder {
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
	}))

	// Delete User by ID
	router.HandleFunc("DELETE /users/{id}", httphandler.Handle(func(r *http.Request) httphandler.Responder {
		id := r.PathValue("id")
		if _, exists := users[id]; !exists {
			return jsonresp.Error(nil, "User not found", http.StatusNotFound)
		}

		delete(users, id)
		return jsonresp.Success[User](nil).WithStatus(http.StatusNoContent)
	}))

	// Download users as CSV
	router.HandleFunc("GET /users/download", httphandler.Handle(func(r *http.Request) httphandler.Responder {
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()

			csvWriter := csv.NewWriter(pw)
			defer csvWriter.Flush()

			for _, it := range users {
				if err := csvWriter.Write([]string{
					it.ID,
					it.Name,
					strconv.Itoa(it.Age),
					it.CreatedAt.Format(time.RFC3339),
				}); err != nil {
					pw.CloseWithError(err)
				}
			}
		}()

		return downloadresp.Attachment(pr, "users.csv")
	}))

	// Redirect
	router.HandleFunc("GET /redirect", httphandler.Handle(func(r *http.Request) httphandler.Responder {
		return httphandler.Redirect("https://google.com", http.StatusSeeOther)
	}))

	// Plain text response
	router.HandleFunc("GET /hello", httphandler.Handle(func(r *http.Request) httphandler.Responder {
		return plainresp.Success("Hello, World!")
	}))

	slog.Info("Server starting on :8080")
	http.ListenAndServe(":8080", router)
}

func validate(input UserInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Age == 0 {
		return fmt.Errorf("age must be > 0")
	}
	return nil
}

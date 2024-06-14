package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

// SignUpPage menampilkan halaman sign up
func SignUpPage(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("views", "signup.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, "Unable to load sign up page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render sign up page", http.StatusInternalServerError)
		return
	}
}

// SignUpHandler memproses sign up dari form HTML atau API JSON
func SignUpHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var username, password string

			if r.Header.Get("Content-Type") == "application/json" {
				// Handle JSON API request
				var request struct {
					Username string `json:"username"`
					Password string `json:"password"`
				}
				err := json.NewDecoder(r.Body).Decode(&request)
				if err != nil {
					http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
					return
				}
				username = request.Username
				password = request.Password
			} else {
				// Handle HTML form request
				err := r.ParseForm()
				if err != nil {
					http.Error(w, "Invalid form submission", http.StatusBadRequest)
					return
				}
				username = r.FormValue("username")
				password = r.FormValue("password")
			}

			// Validate input
			if username == "" || password == "" {
				http.Error(w, "Username dan password harus diisi", http.StatusBadRequest)
				return
			}

			// Hash password before storing (using bcrypt)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Store username and hashedPassword into database (assuming table 'pengguna' exists)
			_, err = db.Exec("INSERT INTO pengguna (username, password) VALUES (?, ?)", username, hashedPassword)
			if err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}

			if r.Header.Get("Content-Type") == "application/json" {
				// Respond with JSON
				response := map[string]string{"message": "User created successfully"}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else {
				// Redirect to login page after successful signup
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// CheckUsernameHandler checks if username already exists
func CheckUsernameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		var exists bool

		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pengguna WHERE username = ?)", username).Scan(&exists)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := map[string]bool{"exists": exists}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

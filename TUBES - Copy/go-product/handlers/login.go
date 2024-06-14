package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

// LoginPage menampilkan halaman login
func LoginPage(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("views", "login.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, "Unable to load login page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render login page", http.StatusInternalServerError)
		return
	}
}

// LoginHandler memproses login dari halaman HTML atau API JSON
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var credentials struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			// Mengurai JSON payload dari body request
			err := json.NewDecoder(r.Body).Decode(&credentials)
			if err != nil {
				http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
				return
			}

			// Query database untuk mendapatkan password terenkripsi
			var hashedPassword string
			err = db.QueryRow("SELECT password FROM pengguna WHERE username = ?", credentials.Username).Scan(&hashedPassword)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Invalid username or password", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Membandingkan password yang dimasukkan dengan yang ada di database
			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
			if err != nil {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			// Jika login berhasil, beri respons sesuai dengan jenis request (HTML atau JSON)
			if r.Header.Get("Content-Type") == "application/json" {
				response := map[string]string{"message": "Login successful"}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else {
				// Redirect ke halaman dashboard setelah login berhasil
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

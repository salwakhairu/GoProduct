package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"
	"text/template"
)

type Product struct {
	Id        int    `json:"id"`
	Nama      string `json:"nama"`
	Harga     string `json:"harga"`
	Deskripsi string `json:"deskripsi"`
}

func NewCreateProductController(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var product Product
			err := json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
				return
			}

			// Insert into database
			result, err := db.Exec("INSERT INTO product (nama, harga, deskripsi) VALUES (?, ?, ?)",
				product.Nama, product.Harga, product.Deskripsi)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}

			// Get the inserted ID
			id, _ := result.LastInsertId()

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Product created successfully",
				"id":      id,
			})

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})

			fp := filepath.Join("views", "create.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// In this case, we would usually pass some initial data to the template,
			// like a struct containing form field values or empty values.
			initialProduct := Product{}
			err = tmpl.Execute(w, initialProduct)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

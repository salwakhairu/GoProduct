package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func NewUpdateProductController(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
			return
		}

		// Extract product ID from URL query parameter
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing product ID"})
			return
		}

		// Convert ID to int
		productID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
			return
		}

		// Execute UPDATE statement
		_, err = db.Exec("UPDATE product SET nama=?, harga=?, deskripsi=? WHERE id=?", product.Nama, product.Harga, product.Deskripsi, productID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Redirect to /product after successful update
		http.Redirect(w, r, "/product", http.StatusSeeOther)
	}
}

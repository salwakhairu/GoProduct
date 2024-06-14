package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// DashboardPage menampilkan halaman dashboard
func DashboardPage(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("views", "dashboard.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

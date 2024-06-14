package routes

import (
	"database/sql"
	"net/http"

	"github.com/salwakhairu/Tugas-Besar-PBW-Sistem-Karyawan/controller"
	"github.com/salwakhairu/Tugas-Besar-PBW-Sistem-Karyawan/handlers"
)

// MapRoutes maps all HTTP routes
func MapRoutes(server *http.ServeMux, db *sql.DB) {
	server.HandleFunc("/", handlers.LoginPage)
	server.HandleFunc("/signup", handlers.SignUpPage)
	server.HandleFunc("/signup/submit", handlers.SignUpHandler(db))
	server.HandleFunc("/login", handlers.LoginPage)
	server.HandleFunc("/login/submit", handlers.LoginHandler(db))
	server.HandleFunc("/dashboard", handlers.DashboardPage)
	server.HandleFunc("/product", controller.NewIndexProduct(db))
	server.HandleFunc("/product/create", controller.NewCreateProductController(db))
	server.HandleFunc("/product/update", controller.NewUpdateProductController(db))
	server.HandleFunc("/product/delete", controller.NewDeleteProductController(db))

	// Routes for API JSON
	server.HandleFunc("/api/login", handlers.LoginHandler(db))
	server.HandleFunc("/api/signup", handlers.SignUpHandler(db))
	server.HandleFunc("/api/check-username", handlers.CheckUsernameHandler(db))
}

package main

import (
	"log"
	"net/http"

	"github.com/salwakhairu/Tugas-Besar-PBW-Sistem-Karyawan/database"
	"github.com/salwakhairu/Tugas-Besar-PBW-Sistem-Karyawan/routes"
)

func main() {
	db := database.InitDatabase()
	defer db.Close()

	server := http.NewServeMux()
	routes.MapRoutes(server, db)

	log.Println("Starting server on :8100")
	err := http.ListenAndServe(":8100", server)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// InitDatabase initializes the database connection and returns the connection pool.
func InitDatabase() *sql.DB {
	// Data Source Name (DSN): username:password@protocol(address)/dbname?param=value
	// Update the DSN string with your actual database connection details.
	dsn := "root:@tcp(localhost:3306)/agen"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}

	// Ping the database to ensure the connection is established
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	return db
}

package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// ConnectDB initializes a MariaDB/MySQL connection and returns a *sql.DB
func ConnectDB() (*sql.DB, error) {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using OS environment variables")
	}

	// Read DB configuration
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	dbname := os.Getenv("DB_NAME")

	if user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("DB_USER, DB_PASSWORD, and DB_NAME must be set")
	}

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		user, password, host, port, dbname,
	)

	// Open DB handle
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %v", err)
	}

	// Optional: Tune connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0)

	// Ping DB to verify actual connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to MariaDB: %v", err)
	}

	log.Println("✅ Successfully connected to MariaDB")
	return db, nil
}

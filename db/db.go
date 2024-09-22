package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func DbConnection() (*sql.DB, error) {
	
	err := godotenv.Load("./.env")
	if err != nil {
		panic("Error loading .env file")
	}

	connectionString := os.Getenv("DB_CONNECTION_STRING")

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, fmt.Errorf("Error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connection Successful")
	return db, nil
}

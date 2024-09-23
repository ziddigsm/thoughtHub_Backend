package db

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func DbConnection() (*gorm.DB, error) {

	err := godotenv.Load("./.env")
	if err != nil {
		panic("Error loading .env file")
	}

	connectionString := os.Getenv("DB_CONNECTION_STRING")

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	return db, nil
}

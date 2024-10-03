package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func DbConnection() (*gorm.DB, error) {

	err := godotenv.Load("./.env")
	if err != nil {
		panic("Error loading .env file")
	}

	connectionString := os.Getenv("DB_CONNECTION_STRING")
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Duration(-1),
		},
	)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	return db, nil
}

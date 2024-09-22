package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		panic("Error loading .env file")
	}
	mgr, err := migrate.New("file://cmd/migrate/migrations", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		fmt.Printf("Migration setup failed: %v", err)
	}
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := mgr.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Migration failed: %v", err)
		}
	}

	if cmd == "down" {
		if err := mgr.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Migration failed: %v", err)
		}
	}

}

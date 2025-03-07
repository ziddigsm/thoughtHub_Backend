package main

import (
	"github.com/ziddigsm/thoughtHub_Backend/cmd/api"
	"github.com/ziddigsm/thoughtHub_Backend/db"
)

func main() {

	db, err := db.DbConnection()
	if err != nil {
		panic("Error initializing database")
	}
	conn, _ := db.DB()
	defer conn.Close()
	server := api.Server("0.0.0.0:8080", db)
	if err := server.Run(); err != nil {
		panic(err)
	}

}

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
	defer db.Close()
	server := api.Server("localhost:8080", nil)
	if err := server.Run(); err != nil {
		panic(err)
	}

}
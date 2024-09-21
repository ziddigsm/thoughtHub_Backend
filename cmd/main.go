package main

import "github.com/ziddigsm/thoughtHub_Backend/cmd/api"

func main() {
	server := api.Server("localhost:8080", nil)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
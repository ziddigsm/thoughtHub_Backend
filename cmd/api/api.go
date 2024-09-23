package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ziddigsm/thoughtHub_Backend/service/users"
	"gorm.io/gorm"
)

type APIServer struct {
	address string
	db      *gorm.DB
}

func Server(address string, db *gorm.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	path := router.PathPrefix("/api/v1").Subrouter()
	userHandler := users.NewHandler(s.db)
	userHandler.InitializeRoutes(path)
	fmt.Println("Server is running on port", s.address)
	return http.ListenAndServe(s.address, router)
}

package api

import (
	"database/sql"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ziddigsm/thoughtHub_Backend/service/users"
)

type APIServer struct {
	address string
	db *sql.DB
}

func Server (address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	path := router.PathPrefix("/api/v1").Subrouter()
	userHandler := users.NewHandler()
	userHandler.InitializeRoutes(path)
	return http.ListenAndServe(s.address, router)
}
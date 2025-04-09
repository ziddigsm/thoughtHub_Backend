package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ziddigsm/thoughtHub_Backend/service/blog"
	"github.com/ziddigsm/thoughtHub_Backend/service/menu"
	"github.com/ziddigsm/thoughtHub_Backend/service/search"
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
	menuHandler := menu.NewHandler(s.db)
	menuHandler.InitializeRoutes(path)
	blogHandler := blog.NewHandler(s.db)
	blogHandler.InitializeRoutes(path)
	searchHandler := search.NewHandler(s.db, blogHandler)
	searchHandler.InitializeRoutes(path)

	enableCors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000", "https://thoughthub.live"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-API-Key"}),
	)
	fmt.Println("Server is running on port", s.address)
	return http.ListenAndServe(s.address, enableCors(router))
}

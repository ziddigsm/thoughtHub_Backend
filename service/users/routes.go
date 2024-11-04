package users

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/get_users", h.GetUsers).Methods("GET")
	router.HandleFunc("/create_user", h.SaveUser).Methods("POST")
	router.HandleFunc("/create_social", h.UpsertSocial).Methods("POST")
}

package menu

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
	router.HandleFunc("/get_menu", h.GetMenu).Methods("GET")
}
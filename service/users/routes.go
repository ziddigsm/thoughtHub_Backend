package users

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.GetUsers).Methods("GET")
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /users API Successfully hit")
}

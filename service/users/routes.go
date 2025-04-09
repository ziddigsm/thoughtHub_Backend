package users

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	    "github.com/ziddigsm/thoughtHub_Backend/utils"

)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) InitializeRoutes(router *mux.Router) {
    router.HandleFunc("/get_users", utils.ApiKeyMiddleware(utils.RateLimitMiddleware(h.GetUsers))).Methods("GET")
    router.HandleFunc("/create_user", utils.ApiKeyMiddleware(utils.RateLimitMiddleware(h.SaveUser))).Methods("POST")
    router.HandleFunc("/create_social", utils.ApiKeyMiddleware(utils.RateLimitMiddleware(h.UpsertSocial))).Methods("POST")
    router.HandleFunc("/save_about", utils.ApiKeyMiddleware(utils.RateLimitMiddleware(h.SaveAbout))).Methods("POST")
    router.HandleFunc("/delete_user", utils.ApiKeyMiddleware(utils.RateLimitMiddleware(h.DeleteUser))).Methods("DELETE")
}

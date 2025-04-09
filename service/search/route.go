package search

import (
	"github.com/gorilla/mux"
	"github.com/ziddigsm/thoughtHub_Backend/service/blog"
	"gorm.io/gorm"
	"github.com/ziddigsm/thoughtHub_Backend/utils"

)

type Handler struct {
	db          *gorm.DB
	blogHandler *blog.Handler
}

func NewHandler(db *gorm.DB, blogHandler *blog.Handler) *Handler {
	return &Handler{
		db:          db,
		blogHandler: blogHandler,
	}
}

func (h *Handler) InitializeRoutes(router *mux.Router) {
    router.HandleFunc("/search_blogs", utils.ApiKeyMiddleware(utils.RateLimitMiddleware(h.SearchBlogs))).Methods("GET")
}

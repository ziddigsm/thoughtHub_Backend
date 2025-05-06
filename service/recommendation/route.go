package recommendation

import (
	"github.com/gorilla/mux"
	"github.com/ziddigsm/thoughtHub_Backend/service/blog"
	"github.com/ziddigsm/thoughtHub_Backend/utils"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	blogHandler *blog.Handler
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db:          db,
		blogHandler: blog.NewHandler(db),
	}
}

func (h *Handler) InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/recommendations", utils.ApiKeyMiddleware(utils.RateLimitMiddleware(h.ValidateRequestBody))).Methods("POST")
}

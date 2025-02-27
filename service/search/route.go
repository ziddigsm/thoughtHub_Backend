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
    router.HandleFunc("/create_blog", h.CreateBlog).Methods("POST")
    router.HandleFunc("/up_likes", h.UpLikes).Methods("GET")
    router.HandleFunc("/get_blogs", h.GetBlogs).Methods("GET")
    router.HandleFunc("/post_comment", h.PostComment).Methods("POST")
    router.HandleFunc("/search_blogs", h.SearchBlogs).Methods("GET")
}
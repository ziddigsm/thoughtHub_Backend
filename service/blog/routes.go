package blog

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
    router.HandleFunc("/create_blog", utils.ApiKeyMiddleware(h.CreateBlog)).Methods("POST")
    router.HandleFunc("/up_likes", utils.ApiKeyMiddleware(h.UpLikes)).Methods("GET")
    router.HandleFunc("/get_blogs", utils.ApiKeyMiddleware(h.GetBlogs)).Methods("GET")
    router.HandleFunc("/post_comment", utils.ApiKeyMiddleware(h.PostComment)).Methods("POST")
    router.HandleFunc("/delete_blog_by_id", utils.ApiKeyMiddleware(h.DeleteBlogByID)).Methods("DELETE")
}

package recommendation

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	service *RecommendationService
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		service: NewService(db),
	}
}

func (h *Handler) InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/recommendations", h.GetSimilarBlogs).Methods("POST")
}

func (h *Handler) GetSimilarBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RecommendationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, `{"error":"No text provided"}`, http.StatusBadRequest)
		return
	}

	if req.UserID == 0 {
		http.Error(w, `{"error":"User ID is required"}`, http.StatusBadRequest)
		return
	}

	response, err := h.service.GetSimilarBlogs(req)
	if err != nil {
		http.Error(w, `{"error":"Failed to get recommendations"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
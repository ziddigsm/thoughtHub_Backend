package users

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/ziddigsm/thoughtHub_Backend/types"
	"github.com/ziddigsm/thoughtHub_Backend/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)


func (h *Handler) SearchBlogs(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    searchTerm := query.Get("search")
    if searchTerm == "" {
        utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("search term is required"))
        return
    }

    var blogs []types.BlogWithName
    var responseblogs []types.DetailedBlog
    response := map[string]interface{}{}
    response["blogs"] = &responseblogs

    // Search in both title and content using ILIKE for case-insensitive search
    if err := h.db.Table("blogs").
        Select("blogs.*, Users.name").
        Joins("LEFT JOIN USERS ON BLOGS.USER_ID = USERS.ID").
        Where("BLOGS.is_active is true AND USERS.is_active is true AND (LOWER(BLOGS.title) LIKE LOWER(?) OR LOWER(BLOGS.content) LIKE LOWER(?))", 
            "%"+searchTerm+"%", 
            "%"+searchTerm+"%").
        Order("BLOGS.created_on desc").
        Find(&blogs).Error; err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to search blogs: %v", err))
        return
    }

    if len(blogs) == 0 {
        response["message"] = "No blogs found matching the search criteria"
        response["totalCount"] = 0
        utils.SuccessResponse(w, http.StatusOK, response)
        return
    }

    // Reuse existing function to get likes and comments
    h.getLikesAndComments(blogs, &responseblogs, w)

    response["totalCount"] = len(blogs)
    response["message"] = "Blogs fetched successfully"
    utils.SuccessResponse(w, http.StatusOK, response)
}
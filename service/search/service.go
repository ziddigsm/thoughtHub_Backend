package search

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ziddigsm/thoughtHub_Backend/types"
	"github.com/ziddigsm/thoughtHub_Backend/utils"
)

func (h *Handler) SearchBlogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchQuery := query.Get("q")
	if searchQuery == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("search query is required"))
		return
	}
	userId, err := strconv.ParseInt(query.Get("user_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("valid user is required"))
		return
	}

	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}

	var blogs []types.BlogWithName
	var responseBlogs []types.DetailedBlog
	response := map[string]interface{}{}
	response["blogs"] = &responseBlogs

	dbQuery := h.db.Table("blogs").Select("blogs.*, Users.name").
		Joins("LEFT JOIN USERS ON BLOGS.USER_ID = USERS.ID").
		Where("BLOGS.is_active is true AND USERS.is_active is true AND (LOWER(BLOGS.title) LIKE LOWER(?) OR LOWER(BLOGS.content) LIKE LOWER(?))",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%").
		Order("BLOGS.created_on desc").
		Limit(int(limit)).
		Offset(int(offset))

	if userId != 0 {
		dbQuery = dbQuery.Where("BLOGS.USER_ID = ?", userId)
	}
	if err := dbQuery.Find(&blogs).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to search blogs: %v", err))
		return
	}

	if len(blogs) == 0 {
		response["message"] = "No blogs found matching your search"
		response["blogs"] = []types.DetailedBlog{}
		utils.SuccessResponse(w, http.StatusOK, response)
		return
	}

	var totalCount int64
	if err := h.db.Table("blogs").
		Joins("LEFT JOIN USERS ON BLOGS.USER_ID = USERS.ID").
		Where("BLOGS.is_active is true AND USERS.is_active is true AND (LOWER(BLOGS.title) LIKE LOWER(?) OR LOWER(BLOGS.content) LIKE LOWER(?))",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%").
		Count(&totalCount).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get total count: %v", err))
		return
	}

	h.blogHandler.GetLikesAndComments(blogs, &responseBlogs, w)

	response["totalCount"] = totalCount
	response["message"] = "Blogs fetched successfully"
	utils.SuccessResponse(w, http.StatusOK, response)
}

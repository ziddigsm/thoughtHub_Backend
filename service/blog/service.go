package blog

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ziddigsm/thoughtHub_Backend/types"
	"github.com/ziddigsm/thoughtHub_Backend/utils"
)

func (h *Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var reqBody types.Blogs
	var likes types.Likes
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse form: %v", err))
		return
	}
	reqBody.UserID, _ = strconv.Atoi(r.FormValue("user_id"))
	reqBody.Title = r.FormValue("title")
	reqBody.Content = r.FormValue("content")
	file, _, err := r.FormFile("blog_image")
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to get file: %v", err))
		return
	}
	defer file.Close()

	reqBody.Blog_image, err = io.ReadAll(file)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to convert file to byte stream: %v", err))
		return
	}
	if err := h.db.Create(&reqBody).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to create blog: %v", err))
		return
	}
	likes.BlogID = reqBody.ID
	if err := h.db.Create(&likes).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to create like table: %v", err))
		return
	}
	response := map[string]interface{}{
		"message": "Blog created successfully",
		"blog":    reqBody,
	}
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *Handler) UpLikes(w http.ResponseWriter, r *http.Request) {
	var reqBody types.Likes
	query := r.URL.Query()
	reqBody.BlogID, _ = strconv.Atoi(query.Get("blog_id"))
	likes, _ := strconv.Atoi(query.Get("likes"))
	reqBody.Likes = likes + 1
	if err := h.db.Model(&reqBody).Where("blog_id = ?", reqBody.BlogID).Update("likes", reqBody.Likes).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to update likes: %v", err))
		return
	}
	response := map[string]interface{}{
		"message": "Likes incremented successfully",
		"likes":   reqBody.Likes,
		"blog_id": reqBody.BlogID,
	}
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *Handler) GetBlogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var blogs []types.BlogWithName
	var responseblogs []types.DetailedBlog
	response := map[string]interface{}{}
	response["blogs"] = &responseblogs
	var blogCount int64
	userId, err := strconv.ParseInt(query.Get("user_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse user_id. Please send a valid user_id: %v", err))
		return
	}
	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse user_id. Please set a valid limit: %v", err))
		return
	}
	if limit != 8 && limit != 9 {
		limit = 9
	}
	offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse user_id. Please set a valid offset: %v", err))
		return
	}
	if offset%limit > 0 {
		offset = offset + (limit - (offset % limit))
	}
	if userId == 0 {
		if err := h.db.Table("blogs").
			Select("blogs.*, Users.name").
			Joins("LEFT JOIN USERS ON BLOGS.USER_ID = USERS.ID").
			Where("BLOGS.is_active is true AND USERS.is_active is true and BLOGS.created_on >= CURRENT_DATE - INTERVAL '3 months' order by BLOGS.created_on desc limit ? offset ? ", limit, offset).
			Find(&blogs).Error; err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get blogs: %v", err))
			return
		}
		if err := h.db.Table("blogs").Select("id").Where("is_active is true and created_on >= CURRENT_DATE - INTERVAL '3 months'").Count(&blogCount).Error; err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get blog count. Please try again: %v", err))
			return
		}
		h.GetLikesAndComments(blogs, &responseblogs, w)
	}
	if userId != 0 {
		if err := h.db.Table("blogs").
			Select("blogs.*, Users.name").
			Joins("LEFT JOIN USERS ON BLOGS.USER_ID = USERS.ID").
			Where("BLOGS.is_active is true AND USERS.is_active is true AND BLOGS.user_id = ? and BLOGS.created_on >= CURRENT_DATE - INTERVAL '3 months' order by BLOGS.created_on desc limit ? offset ? ", userId, limit, offset).
			Find(&blogs).Error; err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get blogs: %v", err))
			return
		}
		if len(blogs) == 0 {
			utils.ErrorResponse(w, http.StatusNoContent, fmt.Errorf("no blogs related to this user"))
			return
		}
		if err := h.db.Table("blogs").Select("id").Where("is_active is true and created_on >= CURRENT_DATE - INTERVAL '3 months' and user_id = ?", userId).Count(&blogCount).Error; err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get blog count. Please try again: %v", err))
			return
		}
		h.GetLikesAndComments(blogs, &responseblogs, w)
	}
	response["totalCount"] = blogCount
	response["message"] = "Blogs fetched successfully"
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *Handler) GetLikesAndComments(blogs []types.BlogWithName, responseblogs *[]types.DetailedBlog, w http.ResponseWriter) {
	var blogIds []int
	for i := range blogs {
		var response types.DetailedBlog
		response.BlogData = blogs[i]
		blogIds = append(blogIds, blogs[i].ID)
		*responseblogs = append(*responseblogs, response)
	}
	var likes []types.Likes
	if err := h.db.Where("blog_id in (?)", blogIds).Find(&likes).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get likes: %v", err))
		return
	}
	likesPerBlogId := make(map[int]int)
	for i := range likes {
		likesPerBlogId[likes[i].BlogID] = likes[i].Likes
	}
	var comments []struct {
		types.Comments
		Name string `json:"name"`
		Mail string `json:"mail"`
	}
	if err := h.db.Table("comments").Select("comments.*, Users.name, Users.mail").Joins("LEFT JOIN USERS ON COMMENTS.USER_ID = USERS.ID").Where("BLOG_ID IN (?) and USERS.is_active is true AND COMMENTS.is_active is true", blogIds).Find(&comments).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get comments: %v", err))
		return
	}
	commentsPerBlogId := make(map[int][]types.DetailedComments)
	for i := range comments {
		var commentsResponse types.DetailedComments
		commentsResponse.Comments = comments[i].Comments
		commentsResponse.Name = comments[i].Name
		commentsResponse.Mail = comments[i].Mail
		commentsPerBlogId[comments[i].BlogID] = append(commentsPerBlogId[comments[i].BlogID], commentsResponse)
	}
	for i := range *responseblogs {
		(*responseblogs)[i].Likes = likesPerBlogId[(*responseblogs)[i].BlogData.ID]
		(*responseblogs)[i].Comments = commentsPerBlogId[(*responseblogs)[i].BlogData.ID]
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var reqBody types.Comments
	err := utils.ParseRequest(r, &reqBody)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse request: %v", err))
		return
	}
	if err := h.db.Create(&reqBody).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to post comment: %v", err))
		return
	}
	response := map[string]interface{}{
		"message": "Comment posted successfully",
		"comment": reqBody,
	}
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *Handler) DeleteBlogByID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var countofBlogs int64
	userId, err := strconv.ParseInt(query.Get("userId"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid userid"))
		return
	}
	blogId, err := strconv.ParseInt(query.Get("blogId"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid blogid"))
		return
	}
	if userId <= 0 || blogId <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid Blog details"))
		return
	}
	if err := h.db.Model(&types.Blogs{}).Where("id = ? and user_id = ? and is_active = true", blogId, userId).Count(&countofBlogs).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get blog count: %v", err))
		return
	}
	if countofBlogs == 0 {
		utils.ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("user unauthorized to delete this blog"))
		return
	}
	if err := h.db.Model(&types.Blogs{}).Where("id = ? AND user_id = ?", blogId, userId).Update("is_active", false).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to delete blog: %v", err))
		return
	}

	response := map[string]interface{}{
		"message": "Blog deleted successfully",
		"blogId":  blogId,
	}
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *Handler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	var reqBody types.Blogs
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse form: %v", err))
		return
	}
	reqBody.UserID, err = strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse user_id: %v", err))
		return
	}
	reqBody.ID, err = strconv.Atoi(r.FormValue("blog_id"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to parse blog_id: %v", err))
		return
	}
	reqBody.Title = r.FormValue("title")
	reqBody.Content = r.FormValue("content")
	file, _, err := r.FormFile("blog_image")
	if err == nil {
		defer file.Close()
		reqBody.Blog_image, err = io.ReadAll(file)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("failed to get file: %v", err))
			return
		}
	}
	fieldsToBeUpdated := map[string]interface{}{}
	if reqBody.Title != "" {
		fieldsToBeUpdated["title"] = reqBody.Title
	}
	if reqBody.Content != "" {
		fieldsToBeUpdated["content"] = reqBody.Content
	}
	if reqBody.Blog_image != nil {
		fieldsToBeUpdated["blog_image"] = reqBody.Blog_image
	}
	if len(fieldsToBeUpdated) == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("no fields to update"))
		return
	}
	fieldsToBeUpdated["updated_on"] = time.Now()
	if err := h.db.Model(&types.Blogs{}).Where("id = ?", reqBody.ID).Updates(fieldsToBeUpdated).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to create blog: %v", err))
		return
	}
	response := map[string]interface{}{
		"message": "Blog created successfully",
		"blog":    reqBody,
	}
	utils.SuccessResponse(w, http.StatusOK, response)
}

package users

import (
	"fmt"
	"net/http"

	"github.com/ziddigsm/thoughtHub_Backend/types"
	"github.com/ziddigsm/thoughtHub_Backend/utils"
	"gorm.io/gorm"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /users API Successfully hit")
}

func (h *Handler) SaveUser(w http.ResponseWriter, r *http.Request) {
	var reqBody types.Users
	if err := utils.ParseRequest(r, &reqBody); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}
	var user types.Users
	if err := h.db.Where("mail = ?", reqBody.Mail).First(&user).Error; err == nil {
		response := map[string] interface{}{
			"message": "User already exists",
            "user_id": user.ID,
		}
		utils.SuccessResponse(w, http.StatusOK , response)
		return
	} else if err != nil  && err != gorm.ErrRecordNotFound {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid emailID: %v", err.Error()))
		return
	}
	if err := h.db.Create(&reqBody).Error; err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to create user: %v", err))
        return
    }

	response := map[string]interface{}{
        "message": "User created successfully",
        "id": reqBody.ID,
		"name": reqBody.Name,
		"mail": reqBody.Mail,
    }
	utils.SuccessResponse(w, http.StatusOK, response)
}
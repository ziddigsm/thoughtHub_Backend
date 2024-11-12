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
    if err := h.db.Where("mail = ? and is_active = ?", reqBody.Mail, true).First(&user).Error; err == nil {
        var socials []types.Socials
        if err := h.db.Where("user_id = ?", user.ID).Find(&socials).Error; err != nil && err != gorm.ErrRecordNotFound {
            utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch social links: %v", err))
            return
        }

        formattedSocials := make(map[string]interface{})
        for _, social := range socials {
            formattedSocials[social.SocialMedia] = social.SocialURL
        }
        response := map[string]interface{}{
            "message":   "User already exists",
            "user_id":   user.ID,
            "name":      user.Name,
            "mail":      user.Mail,
            "username":  user.Username,
            "is_active": user.IsActive,
            "socials":   formattedSocials,
        }
        utils.SuccessResponse(w, http.StatusOK, response)
        return
    } else if err != nil && err != gorm.ErrRecordNotFound {
        utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid emailID: %v", err.Error()))
        return
    }
    if err := h.db.Create(&reqBody).Error; err != nil {
        utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to create user: %v", err))
        return
    }

    response := map[string]interface{}{
        "message": "User created successfully",
        "id":      reqBody.ID,
        "name":    reqBody.Name,
        "mail":    reqBody.Mail,
    }
    utils.SuccessResponse(w, http.StatusOK, response)
}



func (h *Handler) UpsertSocial(w http.ResponseWriter, r *http.Request) {
    var reqBody types.LinkedSocials
    if err := utils.ParseRequest(r, &reqBody); err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
        return
    }
    for _, socialEntry := range reqBody.Socials {
        id, socialMedia, socialURL, err := extractSocialEntry(socialEntry)
        if err != nil {
            utils.ErrorResponse(w, http.StatusBadRequest, err)
            return
        }
        reqData := types.Socials{
            ID:          id,
            UserID:      reqBody.UserID,
            SocialMedia: socialMedia,
            SocialURL:   socialURL,
            IsActive:    reqBody.IsActive,
        }
        if id == 0 {
            if err := h.db.Clauses(clause.OnConflict{
                Columns:   []clause.Column{{Name: "user_id"}, {Name: "social_media"}},
                DoUpdates: clause.AssignmentColumns([]string{"social_media", "social_url", "is_active", "updated_on"}),
            }).Create(&reqData).Error; err != nil {
                utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to insert social link: %v", err))
                return
            }
        } else {
            if err := h.db.Model(&types.Socials{}).Where("id = ?", id).Updates(&reqData).Error; err != nil {
                utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to update social link: %v", err))
                return
            }
        }
    }
    utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
        "user_id": reqBody.UserID,
        "status":  "success",
    })
}

func extractSocialEntry(socialEntry map[string]interface{}) (uint, string, string, error) {
    var id uint
    var socialMedia, socialURL string
    for key, value := range socialEntry {
        if key == "id" && value != nil {
            id = uint(value.(float64))
        } else {
            socialMedia = key
            socialURL, _ = value.(string)
        }
    }
    if socialMedia == "" {
        return 0, "", "", fmt.Errorf("invalid social media format")
    }

    return id, socialMedia, socialURL, nil
}

func (h *Handler) SaveAbout(w http.ResponseWriter, r *http.Request) {
	var reqBody types.Users 
	if err := utils.ParseRequest(r, &reqBody); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}
	 userType := reflect.ValueOf(reqBody)
	 fieldsForUpdate := make(map[string]interface{})
	 res := map[string]interface{}{
		"message": "User updated successfully",
		"user_id": reqBody.ID,
	 }
	i:= 0
	 for key := range userType.NumField() {
		if  (userType.Type().Field(key).Name == "Name" || userType.Type().Field(key).Name == "Username" ){
			fieldsForUpdate[userType.Type().Field(key).Name] = userType.Field(key).Interface()
			i++
			if (userType.Type().Field(key).Name == "Name") {
				res[strings.ToLower(userType.Type().Field(key).Name)] = reqBody.Name
			} else {
				res[strings.ToLower(userType.Type().Field(key).Name)] = reqBody.Username
			}
		}
		if i == 2 {
			break
		}
	 }
	if err := h.db.Model(&reqBody).Where("id = ? and is_active = true", reqBody.ID).Updates(fieldsForUpdate).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to update user: %v", err))
		return
	}
	utils.SuccessResponse(w, http.StatusOK, res)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r * http.Request) {
	query := r.URL.Query()
    id, err := strconv.Atoi(query.Get("id"))
    if err != nil {
        utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid user ID: %v", err))
        return
    }
	if err := h.db.Model(&types.Users{}).Where("id = ? and is_active = true", id).Update("is_active", false).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to delete user: %v", err))
		return
	}
	res := map[string]interface{}{
		"message": "User deleted successfully",
		"user_id": query.Get("id"),
	}
	utils.SuccessResponse(w, http.StatusOK, res)
}
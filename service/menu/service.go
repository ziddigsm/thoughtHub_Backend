package menu

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ziddigsm/thoughtHub_Backend/types"
	"github.com/ziddigsm/thoughtHub_Backend/utils"
)

func (h *Handler) GetMenu(w http.ResponseWriter, r *http.Request) {
	var menu []types.Menu

	query := r.URL.Query()
	isNavbar, _ := strconv.ParseBool(query.Get("is_navbar"))
	if err := h.db.Where("is_navbar = ? and is_active = ?", isNavbar, true).Find(&menu).Error; err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch menu: %v", err))
		return
	}
	utils.SuccessResponse(w, http.StatusOK, menu)
}
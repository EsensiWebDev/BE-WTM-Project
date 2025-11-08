package user_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateSetting updates user settings
// @Summary Update user settings
// @Description Update user settings
// @Tags Profile
// @Produce json
// @Router /profile/setting [put]
// @Param request body userdto.UpdateSettingRequest true "Update Setting Request"
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully updated user settings"
func (uh *UserHandler) UpdateSetting(c *gin.Context) {
	ctx := c.Request.Context()

	// Bind the request data to the UpdateSettingRequest struct
	var req userdto.UpdateSettingRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Error validating request:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}

		// fallback: unknown validation error
		logger.Error(ctx, "Unexpected validation error", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := uh.userUsecase.UpdateSetting(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating user settings:", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update user settings: %s", err.Error()))
		return
	}

	utils.ClearRefreshCookie(c, uh.config.URL, uh.config.SecureService)

	response.Success(c, nil, "Successfully updated user settings")
}

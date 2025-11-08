package notification_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/notifdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateNotificationSettings godoc
// @Summary Update notification settings
// @Description Update user's notification settings
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body notifdto.UpdateNotificationSettingRequest true "Update Notification Settings Request"
// @Router /notifications/settings [put]
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully updated notification settings"
func (nh *NotificationHandler) UpdateNotificationSettings(c *gin.Context) {
	ctx := c.Request.Context()

	var req notifdto.UpdateNotificationSettingRequest
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

		logger.Error(ctx, "Unexpected validation error", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := nh.notifUsecase.UpdateNotificationSetting(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating notification settings:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to update notification settings")
		return
	}

	response.Success(c, nil, "Successfully updated notification settings")
}

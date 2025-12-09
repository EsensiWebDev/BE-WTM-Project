package notification_handler

import (
	"net/http"
	"wtm-backend/internal/dto/notifdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ReadNotification godoc
// @Summary Read Notification
// @Description Read Notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body notifdto.ReadNotificationRequest true "Read Notification Request"
// @Router /notifications/read [put]
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully read notification"
func (nh *NotificationHandler) ReadNotification(c *gin.Context) {
	ctx := c.Request.Context()

	var req notifdto.ReadNotificationRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind request payload", err)
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := nh.notifUsecase.ReadNotification(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to read notification", err)
		response.Error(c, http.StatusInternalServerError, "Failed to read notification")
		return
	}

	response.Success(c, nil, "Successfully read notification")
}

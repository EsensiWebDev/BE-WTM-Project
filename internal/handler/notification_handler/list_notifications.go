package notification_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/notifdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListNotifications godoc
// @Summary List Notifications
// @Description Get a list of notifications with pagination
// @Tags Notifications
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page"
// @Success 200 {object} response.ResponseWithPagination{data=[]entity.Notification} "Successfully retrieved notifications"
// @Security BearerAuth
// @Router /notifications [get]
func (nh *NotificationHandler) ListNotifications(c *gin.Context) {
	ctx := c.Request.Context()

	var req notifdto.ListNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := nh.notifUsecase.ListNotifications(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error listing notifications:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to list notifications")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved notifications"

	var hotels []entity.Notification
	if resp != nil {
		hotels = resp.Notifications
		if len(resp.Notifications) == 0 {
			message = "No notifications found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, hotels, message, pagination)
}

package email_handler

import (
	"net/http"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ListEmailLogs godoc
// @Summary List Email Logs
// @Description Retrieve a paginated list of email logs using query parameters.
// @Tags Email
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithPagination{data=[]emaildto.EmailLogResponse} "Successfully retrieved list of email logs"
// @Router /email/logs [get]
func (eh *EmailHandler) ListEmailLogs(c *gin.Context) {
	ctx := c.Request.Context()

	var req emaildto.ListEmailLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := eh.emailUsecase.ListEmailLogs(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching email logs:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of email logs")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of email logs"

	var emailLogs []emaildto.EmailLogResponse
	if resp != nil {
		emailLogs = resp.EmailLogs
		if len(resp.EmailLogs) == 0 {
			message = "No email logs found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, emailLogs, message, pagination)
}

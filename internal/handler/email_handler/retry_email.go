package email_handler

import (
	"net/http"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RetryEmail godoc
// @Summary Retry Failed Email
// @Description Retry sending a failed email
// @Tags Email
// @Accept json
// @Produce json
// @Param request body emaildto.RetryEmailRequest true "Retry Email Request"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=emaildto.RetryEmailResponse} "Email retry initiated"
// @Router /email/logs/retry [post]
func (eh *EmailHandler) RetryEmail(c *gin.Context) {
	ctx := c.Request.Context()

	var req emaildto.RetryEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := eh.emailUsecase.RetryEmail(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error retrying email:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to retry email")
		return
	}

	if !resp.Success {
		response.Error(c, http.StatusBadRequest, resp.Message)
		return
	}

	response.Success(c, resp, resp.Message)
}


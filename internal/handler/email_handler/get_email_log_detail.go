package email_handler

import (
	"net/http"
	"strconv"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GetEmailLogDetail godoc
// @Summary Get Email Log Detail
// @Description Retrieve detailed information about a specific email log
// @Tags Email
// @Accept json
// @Produce json
// @Param id path int true "Email Log ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=emaildto.EmailLogDetailResponse} "Successfully retrieved email log detail"
// @Router /email/logs/{id} [get]
func (eh *EmailHandler) GetEmailLogDetail(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error(ctx, "Invalid email log ID:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid email log ID")
		return
	}

	resp, err := eh.emailUsecase.GetEmailLogDetail(ctx, uint(id))
	if err != nil {
		logger.Error(ctx, "Error fetching email log detail:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get email log detail")
		return
	}

	response.Success(c, resp, "Successfully retrieved email log detail")
}

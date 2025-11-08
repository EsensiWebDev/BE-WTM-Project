package email_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// EmailTemplate godoc
// @Summary Get Email Templates
// @Description Retrieve a list of email templates
// @Tags Email
// @Produce json
// @Success 200 {object} response.Response{data=[]emaildto.EmailTemplateResponse} "Successfully retrieved email templates"
// @Router /email/template [get]
// @Security BearerAuth
func (eh *EmailHandler) EmailTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	templates, err := eh.emailUsecase.EmailTemplate(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting email templates:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get email templates")
		return
	}

	if templates == nil {
		logger.Error(ctx, "No email template found")
		response.Success(c, nil, "No email template found")
		return
	}

	response.Success(c, templates, "Successfully retrieved email templates")
}

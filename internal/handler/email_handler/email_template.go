package email_handler

import (
	"net/http"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// EmailTemplate godoc
// @Summary Get Email Templates
// @Description Retrieve a list of email templates
// @Tags Email
// @Produce json
// @Param type query string false "Type of email template (option: 'confirm', 'cancel')"
// @Success 200 {object} response.Response{data=[]emaildto.EmailTemplateResponse} "Successfully retrieved email templates"
// @Router /email/template [get]
// @Security BearerAuth
func (eh *EmailHandler) EmailTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	var req emaildto.EmailTemplateRequest
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

	templates, err := eh.emailUsecase.EmailTemplate(ctx, &req)
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

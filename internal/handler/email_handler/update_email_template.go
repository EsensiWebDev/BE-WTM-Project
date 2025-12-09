package email_handler

import (
	"net/http"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UpdateEmailTemplate godoc
// @Summary Update Email Template
// @Description Update an existing email template
// @Tags Email
// @Accept multipart/form-data
// @Produce json
// @Param type formData string true "Template Type"
// @Param subject formData string true "Template Subject"
// @Param body formData string true "Template Body dalam html format"
// @Param signature_text formData string false "Template Signature Text dalam html format"
// @Param signature_image formData file false "Template Signature Image"
// @Success 200 {object} response.Response "Successfully updated email template"
// @Security BearerAuth
// @Router /email/template [put]
func (eh *EmailHandler) UpdateEmailTemplate(c *gin.Context) {
	ctx := c.Request.Context()

	var req emaildto.UpdateEmailTemplateRequest
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

	if err := eh.emailUsecase.UpdateEmailTemplate(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating email template:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed updating email template")
		return
	}

	response.Success(c, nil, "Successfully updated email template")
}

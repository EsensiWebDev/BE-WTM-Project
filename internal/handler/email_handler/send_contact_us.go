package email_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// SendContactUs godoc
// @Summary Send Contact Us Email
// @Description Send a contact us email
// @Tags Email
// @Accept json
// @Produce json
// @Param request body emaildto.SendContactUsEmailRequest true "Contact Us Email Request"
// @Success 200 {object} response.Response "Successfully sent contact us email"
// @Router /email/contact-us [post]
func (eh *EmailHandler) SendContactUs(c *gin.Context) {
	ctx := c.Request.Context()

	var req emaildto.SendContactUsEmailRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := eh.emailUsecase.SendContactUsEmail(ctx, &req); err != nil {
		logger.Error(ctx, "Error sending contact us email:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to send contact us email")
		return
	}

	response.Success(c, nil, "Successfully sent contact us email")

}

package email_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// SendContactUs godoc
// @Summary Send Contact Us Email
// @Description Send a contact us email
// @Tags Email
// @Accept json
// @Produce json
// @Param name query string true "Name"
// @Param email query string true "Email"
// @Param subject query string true "Subject"
// @Param department query string true "Department"
// @Param type query string true "Type" Enums(general, booking)
// @Param booking_code query string false "Booking Code"
// @Param sub_booking_code query string false "Sub Booking Code"
// @Param message query string true "Message"
// @Success 200 {object} response.Response "Successfully sent contact us email"
// @Router /email/contact-us [post]
func (eh *EmailHandler) SendContactUs(c *gin.Context) {
	ctx := c.Request.Context()

	var req emaildto.SendContactUsEmailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := eh.emailUsecase.SendContactUsEmail(ctx, &req); err != nil {
		logger.Error(ctx, "Error sending contact us email:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to send contact us email")
		return
	}

	response.Success(c, nil, "Successfully sent contact us email")

}

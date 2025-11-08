package auth_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ForgotPassword godoc
// @Summary      Forgot Password
// @Description  Initiates the password reset process by sending a reset link to the user's email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      authdto.ForgotPasswordRequest  true  "Forgot Password Payload"
// @Failure      200      {object}  response.Response  "Password reset link has been sent to your email if it exists in our system"
// @Success      200      {object}  response.Response{data=authdto.ForgotPasswordResponse}  "Password reset link has been sent to your email"
// @Router       /forgot-password [post]
func (ah *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req authdto.ForgotPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate the request data
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

	resp, err := ah.authUsecase.ForgotPassword(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error in ForgotPassword usecase:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to process forgot password request")
		return
	}

	if resp == nil {
		// In case of non-existing email, we do not reveal that information
		response.Success(c, nil, "Password reset link has been sent to your email if it exists in our system")
		return
	}

	response.Success(c, resp, "Password reset link has been sent to your email")
}

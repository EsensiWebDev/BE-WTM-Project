package auth_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ResetPassword godoc
// @Summary      Reset Password
// @Description  Resets the user's password using a valid reset token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      authdto.ResetPasswordRequest  true  "Reset Password Payload"
// @Success      200      {object}  response.Response  "Password has been reset
// @Router       /reset-password [post]
func (ah *AuthHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req authdto.ResetPasswordRequest
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

	if err := ah.authUsecase.ResetPassword(ctx, &req); err != nil {
		logger.Error(ctx, "Error in ResetPassword usecase:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	response.Success(c, nil, "Password has been reset successfully")
}

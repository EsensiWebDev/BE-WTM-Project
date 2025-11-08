package auth_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ValidateTokenResetPassword godoc
// @Summary      Validate Reset Password Token
// @Description  Validates the password reset token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      authdto.ValidateTokenResetPasswordRequest true  "Validate Reset Password Token Payload"
// @Success      200      {object}  response.Response{data=authdto.ValidateTokenResetPasswordResponse}  "Reset password token is valid"
// @Router       /reset-password [get]
func (ah *AuthHandler) ValidateTokenResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req authdto.ValidateTokenResetPasswordRequest
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

	resp, err := ah.authUsecase.ValidateTokenResetPassword(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error in ValidateTokenResetPassword usecase:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to validate reset password token")
		return
	}

	response.Success(c, resp, "Reset password token is valid")
}

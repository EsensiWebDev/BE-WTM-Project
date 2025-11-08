package auth_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// Login logs in a user with the provided credentials.
// @Summary User Login
// @Description Logs in a user with username and password
// @Tags Auth
// @Produce json
// @Param request body authdto.LoginRequest true "Login Request"
// @Success 200 {object} response.ResponseWithData{data=authdto.LoginResponse} "Login successful"
// @Router /login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var req authdto.LoginRequest

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

	respLogin, refreshToken, err := ah.authUsecase.Login(ctx, &req)
	if err != nil {
		logger.Warn(ctx, "Login failed:", err.Error())
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SetRefreshCookie(c, refreshToken, ah.config.URL, int(ah.config.DurationRefreshToken.Seconds()), ah.config.SecureService)

	response.Success(c, respLogin, "Login successful")
}

package auth_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// Logout logs out the user by clearing the refresh token cookie.
// @Summary User Logout
// @Description Logs out the user by clearing the refresh token cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Response "Logout successfully"
// @Router /logout [post]
// @Security BearerAuth
func (ah *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	_, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Error(ctx, "Error retrieving refresh token cookie:", err.Error())
		response.Error(c, http.StatusUnauthorized, "Cookie failed")
		return
	}

	if err := ah.authUsecase.Logout(ctx); err != nil {
		logger.Error(ctx, "Error during logout:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Logout failed")
		return
	}

	utils.ClearRefreshCookie(c, ah.config.URL, ah.config.SecureService)

	response.Success(c, nil, "Logout successfully")
}

package auth_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RefreshToken refreshes the user's access token using the refresh token stored in a cookie.
// @Summary Refresh User Token
// @Description Refreshes the user's access token using the refresh token stored in a cookie
// @Tags Auth
// @Produce json
// @Success 200 {object} response.ResponseWithData{data=authdto.LoginResponse} "Token refreshed successfully"
// @Router /refresh-token [get]
func (ah *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Error(ctx, "Error retrieving refresh token from cookie:", err.Error())
		response.Error(c, http.StatusUnauthorized, "Cookie failed")
		return
	}

	resp, err := ah.authUsecase.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		logger.Error(ctx, "Error refreshing token:", err.Error())
		utils.ClearRefreshCookie(c, ah.config.URL, ah.config.SecureService)
		response.Error(c, http.StatusUnauthorized, "Cookie failed")
		return
	}

	response.Success(c, resp, "Token refreshed successfully")
}

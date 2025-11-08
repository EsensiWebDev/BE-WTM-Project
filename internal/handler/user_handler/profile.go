package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// Profile godoc
// @Summary Get authenticated user's profile
// @Description Retrieve the profile information of the currently authenticated user based on the JWT token provided in the Authorization header.
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithData{data=userdto.ProfileResponse} "Successfully retrieved user profile"
// @Router /profile [get]
func (uh *UserHandler) Profile(c *gin.Context) {
	ctx := c.Request.Context()

	dataUser, err := uh.userUsecase.Profile(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting user profile:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	if dataUser == nil {
		logger.Error(ctx, "User not found")
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, dataUser, "Successfully retrieved user profile")
}

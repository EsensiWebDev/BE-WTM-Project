package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListRoleAccess godoc
// @Summary Get list of role access
// @Description Retrieve a list of available role access configurations.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseWithData{data=[]userdto.ListRoleAccessResponse} "Successfully retrieved role access list"
// @Security BearerAuth
// @Router /role-access [get]
func (uh *UserHandler) ListRoleAccess(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := uh.userUsecase.ListRoleAccess(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting role access", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to fetch role access")
		return
	}

	message := "Successfully retrieved role access list"

	if len(data) == 0 {
		message = "No role access found"
	}

	response.Success(c, data, message)
}

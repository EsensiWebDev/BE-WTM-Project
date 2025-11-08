package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateStatusUser godoc
// @Summary Update user status
// @Description Update the status of a user (active/inactive).
// @Tags User
// @Accept json
// @Produce json
// @Param request body userdto.UpdateStatusUserRequest true "Status of the user"
// @Success 200 {object} response.Response "Successfully updated user settings"
// @Security BearerAuth
// @Router /users/status [post]
func (uh *UserHandler) UpdateStatusUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.UpdateStatusUserRequest
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

		logger.Error(ctx, "Unexpected validation error", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := uh.userUsecase.UpdateStatusUser(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating user status", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to update user status")
		return
	}

	response.Success(c, nil, "Successfully updated user settings")
}

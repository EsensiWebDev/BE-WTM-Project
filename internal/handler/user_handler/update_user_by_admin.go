package user_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateUserByAdmin godoc
// @Summary Update a user by admin
// @Description Update a user with the provided details
// @Tags User
// @Produce json
// @Param request body userdto.UpdateUserByAdminRequest true "Update user request payload"
// @Router /users [put]
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully updated user"
func (uh *UserHandler) UpdateUserByAdmin(c *gin.Context) {
	var req userdto.UpdateUserByAdminRequest
	ctx := c.Request.Context()

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

	if err := uh.userUsecase.UpdateUserByAdmin(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating user", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %s", err.Error()))
		return
	}

	response.Success(c, nil, "Successfully updated user")

}

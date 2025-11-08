package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateRoleAccess godoc
// @Summary Update role access configuration
// @Description Update the access permissions associated with a specific role by providing the updated configuration in JSON format.
// @Tags User
// @Accept json
// @Produce json
// @Param request body userdto.UpdateRoleAccessRequest true "Payload to update role access"
// @Success 200 {object} response.Response "Successfully updated role access"
// @Security BearerAuth
// @Router /role-access [put]
func (uh *UserHandler) UpdateRoleAccess(c *gin.Context) {
	ctx := c.Request.Context()
	var req userdto.UpdateRoleAccessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		if ve := utils.ParseValidationErrors(err); ve != nil {
			logger.Error(ctx, "Error validating request:", err.Error())
			response.ValidationError(c, ve)
			return
		}

		logger.Error(ctx, "Error validating request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := uh.userUsecase.UpdateRoleAccess(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating role access:", err.Error())
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil, "Successfully updated role access")
}

package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListUsersByRole godoc
// @Summary Get users by role
// @Description Retrieve a paginated list of users filtered by role using path and query parameters.
// @Tags User
// @Accept json
// @Produce json
// @Param role path string true "Role name (e.g. admin, agent, support, super_admin)"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter users"
// @Success 200 {object} response.ResponseWithPagination "Successfully retrieved users by role"
// @Security BearerAuth
// @Router /users/by-role/{role} [get]

func (uh *UserHandler) ListUsersByRole(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.ListUsersByRoleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	roleParam := c.Param("role")
	req.Role = roleParam

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

	resp, total, err := uh.userUsecase.ListUsersByRole(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error getting agent companies", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get users by role")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Users by role retrieved successfully"

	var users []userdto.ListUsersByRoleData
	if resp != nil {
		users = resp.User
		if len(users) == 0 {
			message = "No user by role found"
		}
	}

	response.SuccessWithPagination(c, users, message, pagination)
}

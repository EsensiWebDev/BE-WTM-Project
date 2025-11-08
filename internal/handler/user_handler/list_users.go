package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListUsers godoc
// @Summary Get users
// @Description Retrieve a paginated list of users using query parameters.
// @Tags User
// @Accept json
// @Produce json
// @Param role query string false "Role name to filter users (e.g. admin, agent, support, super_admin)"
// @Param agent_company_id query int false "Id of the agent company to filter users"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter users"
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithPagination{data=[]userdto.ListUserData} "Successfully retrieved users"
// @Router /users [get]
func (uh *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.ListUsersRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	req.Scope = constant.ScopeManagement

	resp, err := uh.userUsecase.ListUsers(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching users:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get users")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved users"

	if resp == nil || resp.Users == nil || len(resp.Users) == 0 {
		message = "No users found"
		response.EmptyList(c, message, pagination)
		return
	}

	users := resp.Users
	pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))

	response.SuccessWithPagination(c, users, message, pagination)
}

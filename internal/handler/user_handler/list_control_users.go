package user_handler

import (
	"net/http"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ListControlUsers godoc
// @Summary      List Control Users
// @Description  Retrieve a paginated list of control users using query parameters.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        page               query     int     false  "Page number for pagination"
// @Param        limit              query     int     false  "Number of items per page"
// @Param        search             query     string  false  "Search keyword to filter users"
// @Security     BearerAuth
// @Success      200  {object}  response.ResponseWithPagination{data=[]userdto.ListUserData} "Successfully retrieved users"
// @Router       /users/control [get]
func (uh *UserHandler) ListControlUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.ListUsersRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Role == "" {
		req.Role = constant.RoleAgent
	}
	req.Scope = constant.ScopeControl

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

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

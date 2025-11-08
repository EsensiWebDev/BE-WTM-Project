package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListStatusUsers godoc
// @Summary List status users
// @Description Retrieve a list of status users with pagination
// @Tags User
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page"
// @Param page query int false "Page number for pagination"
// @Param search query string false "Search keyword to filter users"
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithPagination{data=[]entity.StatusUser} "Successfully retrieved status users"
// @Router /users/status [get]

func (uh *UserHandler) ListStatusUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.ListStatusUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, total, err := uh.userUsecase.ListStatusUsers(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error listing users by status:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to list status users")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved status users"

	var statusUsers []entity.StatusUser
	if resp != nil {
		statusUsers = resp.StatusUsers
		if len(resp.StatusUsers) == 0 {
			message = "No status users found"
		}
	}

	response.SuccessWithPagination(c, statusUsers, message, pagination)
}

package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListUsersByAgentCompany godoc
// @Summary Get users by agent company
// @Description Retrieve a paginated list of users associated with a specific agent company using query parameters.
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "Id of the agent company"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter users"
// @Success 200 {object} response.ResponseWithPagination{data=[]userdto.ListUsersByAgentCompanyData} "Successfully retrieved users by agent company"
// @Security BearerAuth
// @Router /users/by-agent-company/{id} [get]
func (uh *UserHandler) ListUsersByAgentCompany(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.ListUsersByAgentCompanyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logger.Error(ctx, "Invalid agent company Id in path:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid agent company Id")
		return
	}
	req.ID = uint(id)

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

	resp, total, err := uh.userUsecase.ListUsersByAgentCompany(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error getting users by agent company:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get users by agent company")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved users by agent company"

	var users []userdto.ListUsersByAgentCompanyData
	if resp != nil {
		users = resp.Users
		if len(users) == 0 {
			message = "No users found for this agent company"
		}
	}

	response.SuccessWithPagination(c, users, message, pagination)
}

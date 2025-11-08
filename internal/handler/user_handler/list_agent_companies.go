package user_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListAgentCompanies godoc
// @Summary Get list of agent companies
// @Description Retrieve a paginated list of agent companies using query parameters.
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter agent companies"
// @Success 200 {object} response.ResponseWithPagination{data=[]entity.AgentCompany} "Successfully retrieved agent companies"
// @Security BearerAuth
// @Router /users/agent-companies [get]
func (uh *UserHandler) ListAgentCompanies(c *gin.Context) {
	ctx := c.Request.Context()

	var req userdto.ListAgentCompaniesRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, total, err := uh.userUsecase.ListAgentCompanies(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error getting agent companies", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get agent companies")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved agent companies"

	var agentCompanies []entity.AgentCompany
	if resp != nil {
		agentCompanies = resp.AgentCompanies
		if len(agentCompanies) == 0 {
			message = "No agent companies found"
		}
	}

	response.SuccessWithPagination(c, agentCompanies, message, pagination)
}

package promo_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/response"
)

// ListPromoForAgent godoc
// @Summary List promos for agent
// @Description List promos for agent
// @Tags Promo
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter promos by name,description, and code"
// @Success 200 {object} response.ResponseWithPagination{data=[]promodto.PromosForAgent} "Successfully retrieved list of promos"
// @Security BearerAuth
// @Router /promos/agent [get]
func (ph *PromoHandler) ListPromoForAgent(c *gin.Context) {
	ctx := c.Request.Context()

	var req promodto.ListPromosForAgentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := ph.promoUsecase.ListPromosForAgent(ctx, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get list of promos")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of promos"

	var promos []promodto.PromosForAgent
	if resp != nil {
		promos = resp.Data
		if len(promos) == 0 {
			message = "No promos found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, promos, message, pagination)

}

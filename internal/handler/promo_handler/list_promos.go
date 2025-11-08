package promo_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/response"
)

// ListPromos godoc
// @Summary List Promos
// @Description Retrieve a list of promos with pagination.
// @Tags Promo
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter promos by name"
// @Success 200 {object} response.ResponseWithPagination{data=[]promodto.PromoResponse} "Successfully retrieved list of promos"
// @Security BearerAuth
// @Router /promos [get]
func (ph *PromoHandler) ListPromos(c *gin.Context) {
	ctx := c.Request.Context()

	var req promodto.ListPromosRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, total, err := ph.promoUsecase.ListPromos(ctx, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get list of promos")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved list of promos"

	var promos []promodto.PromoResponse
	if resp != nil {
		promos = resp.Promos
		if len(resp.Promos) == 0 {
			message = "No promos found"
		}
	}

	response.SuccessWithPagination(c, promos, message, pagination)
}

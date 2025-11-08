package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListUnassignedPromos godoc
// @Summary List Unassigned Promos
// @Description List promos that are not assigned to any promo group.
// @Tags Promo Group
// @Accept       json
// @Produce      json
// @Param        page  query  int  false  "Page number for pagination"
// @Param        limit query  int  false  "Number of items per page"
// @Param        search query  string  false  "Search term for filtering promos"
// @Param		 promo_group_id query int true "Promo Group ID to exclude promos already assigned to this group"
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithPagination{data=[]promogroupdto.ListUnassignedPromoData} "Successfully retrieved unassigned promos"
// @Router /promo-groups/unassigned-promos [get]
func (pgh *PromoGroupHandler) ListUnassignedPromos(c *gin.Context) {
	ctx := c.Request.Context()

	var req promogroupdto.ListUnassignedPromosRequest
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

	resp, err := pgh.promoGroupUsecase.ListUnassignedPromos(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error listing unassigned promos:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to list unassigned promos")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved unassigned promos"

	if resp == nil || resp.Promos == nil || len(resp.Promos) == 0 {
		message = "No unassigned promos found"
		response.EmptyList(c, message, pagination)
		return
	}

	promos := resp.Promos
	pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))

	response.SuccessWithPagination(c, promos, message, pagination)

}

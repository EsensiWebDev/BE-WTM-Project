package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListPromoGroupPromos godoc
// @Summary Get detailed promos in a promo group
// @Description Retrieve a paginated list of promos that belong to a specific promo group using query parameters.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param id query int true "Promo Group ID"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithPagination{data=[]promogroupdto.ListPromoGroupPromosData} "Successfully retrieved promo group details"
// @Router /promo-groups/promos [get]
func (pgh *PromoGroupHandler) ListPromoGroupPromos(c *gin.Context) {
	ctx := c.Request.Context()

	var req promogroupdto.ListPromoGroupPromosRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusUnprocessableEntity, "Invalid request payload")
		return
	}

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

	if req.Page < 1 {
		req.Page = 1
	}

	resp, total, err := pgh.promoGroupUsecase.ListPromoGroupPromos(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error getting detail promo group:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get detail promo group")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved promo group details"

	var promos []promogroupdto.ListPromoGroupPromosData
	if resp != nil {
		promos = resp.Promos

		if len(promos) == 0 {
			message = "No detail promo group found"
		}
	}

	response.SuccessWithPagination(c, promos, message, pagination)
}

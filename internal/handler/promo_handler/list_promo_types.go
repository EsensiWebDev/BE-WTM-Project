package promo_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListPromoTypes godoc
// @Summary List Promo Types
// @Description Retrieve a list of promo types with pagination.
// @Tags Promo
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter promo types by name"
// @Success 200 {object} response.ResponseWithPagination{data=[]entity.PromoType} "Successfully retrieved list of promo types"
// @Security BearerAuth
// @Router /promos/types [get]
func (ph *PromoHandler) ListPromoTypes(c *gin.Context) {
	ctx := c.Request.Context()

	var req promodto.ListPromoTypesRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	resp, total, err := ph.promoUsecase.ListPromoTypes(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching promo types:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of promo types")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved list of promo types"

	var promoTypes []entity.PromoType
	if resp != nil {
		promoTypes = resp.PromoTypes
		if len(resp.PromoTypes) == 0 {
			message = "No facilities found"
		}
	}

	response.SuccessWithPagination(c, promoTypes, message, pagination)
}

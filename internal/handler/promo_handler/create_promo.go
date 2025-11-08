package promo_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// CreatePromo godoc
// @Summary Create Promo
// @Description Create a new promo with the specified details.
// @Tags Promo
// @Accept json
// @Produce json
// @Param request body promodto.UpsertPromoRequest true "Promo details"
// @Success 200 {object} response.Response "Successfully created promo"
// @Security BearerAuth
// @Router /promos [post]
func (ph *PromoHandler) CreatePromo(c *gin.Context) {
	ctx := c.Request.Context()

	var req *promodto.UpsertPromoRequest
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

	err := ph.promoUsecase.UpsertPromo(ctx, req, nil)
	if err != nil {
		logger.Error(ctx, "Error creating promo:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error creating promo")
		return
	}

	response.Success(c, nil, "Successfully created promo")
}

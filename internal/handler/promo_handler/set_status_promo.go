package promo_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// SetStatusPromo godoc
// @Summary      Set Status Promo
// @Description  API to set the status of a promo (active/inactive)
// @Tags         Promo
// @Accept       multipart/form-data
// @Produce      json
// @Param        promo_id   formData  string  true  "Promo ID"
// @Param        is_active  formData  bool    true  "IsActive to set (true for active, false for inactive)"
// @Success      200      {object}  response.Response "Successfully set promo status"
// @Security     BearerAuth
// @Router       /promos/status [put]
func (ph *PromoHandler) SetStatusPromo(c *gin.Context) {
	ctx := c.Request.Context()

	var req promodto.SetStatusPromoRequest
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

	err := ph.promoUsecase.SetStatusPromo(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error setting promo status:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error setting promo status")
		return
	}

	response.Success(c, nil, "Successfully set promo status")
}

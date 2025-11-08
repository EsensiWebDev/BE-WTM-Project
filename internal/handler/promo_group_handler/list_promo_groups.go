package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListPromoGroups godoc
// @Summary Get list of promo groups
// @Description Retrieve a paginated list of promo groups based on optional filters such as search keyword, page, and limit.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter promo groups"
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithPagination{data=[]entity.PromoGroup} "Successfully retrieved list of promo groups"
// @Router /promo-groups [get]
func (pgh *PromoGroupHandler) ListPromoGroups(c *gin.Context) {
	ctx := c.Request.Context()

	var req promogroupdto.ListPromoGroupRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusUnprocessableEntity, "Invalid request payload")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}

	resp, total, err := pgh.promoGroupUsecase.ListPromoGroups(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error getting promo groups", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get promo groups")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved list of promo groups"

	var promoGroups []entity.PromoGroup
	if resp != nil {
		promoGroups = resp.PromoGroups
		if len(promoGroups) == 0 {
			message = "No promo groups found"
		}
	}

	response.SuccessWithPagination(c, promoGroups, message, pagination)
}

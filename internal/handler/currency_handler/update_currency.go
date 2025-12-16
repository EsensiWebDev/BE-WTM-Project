package currency_handler

import (
	"net/http"
	"strconv"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/currencydto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UpdateCurrency godoc
// @Summary Update Currency
// @Description Update an existing currency (code cannot be changed)
// @Tags Currency
// @Accept json
// @Produce json
// @Param id path int true "Currency ID"
// @Param request body currencydto.UpdateCurrencyRequest true "Update Currency Request"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=currencydto.CurrencyResponse} "Successfully updated currency"
// @Router /currencies/{id} [put]
func (ch *CurrencyHandler) UpdateCurrency(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error(ctx, "Invalid currency ID", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid currency ID")
		return
	}

	var req currencydto.UpdateCurrencyRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	currency := &entity.Currency{
		ID:       uint(id),
		Name:     req.Name,
		Symbol:   req.Symbol,
		IsActive: req.IsActive,
	}

	updated, err := ch.currencyUsecase.UpdateCurrency(ctx, currency)
	if err != nil {
		logger.Error(ctx, "Error updating currency", err.Error())
		if utils.ParseValidationErrors(err) != nil {
			response.ValidationError(c, utils.ParseValidationErrors(err))
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, currencydto.ToCurrencyResponse(updated), "Successfully updated currency")
}

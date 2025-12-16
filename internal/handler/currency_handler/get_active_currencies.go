package currency_handler

import (
	"net/http"
	"wtm-backend/internal/dto/currencydto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GetActiveCurrencies godoc
// @Summary Get Active Currencies
// @Description Retrieve only active currencies
// @Tags Currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]currencydto.CurrencyResponse} "Successfully retrieved active currencies"
// @Router /currencies/active [get]
func (ch *CurrencyHandler) GetActiveCurrencies(c *gin.Context) {
	ctx := c.Request.Context()

	currencies, err := ch.currencyUsecase.GetActiveCurrencies(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting active currencies", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get active currencies")
		return
	}

	var currencyResponses []currencydto.CurrencyResponse
	for _, currency := range currencies {
		currencyResponses = append(currencyResponses, currencydto.ToCurrencyResponse(&currency))
	}

	response.Success(c, currencyResponses, "Successfully retrieved active currencies")
}

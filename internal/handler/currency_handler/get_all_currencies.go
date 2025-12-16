package currency_handler

import (
	"net/http"
	"wtm-backend/internal/dto/currencydto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GetAllCurrencies godoc
// @Summary Get All Currencies
// @Description Retrieve all currencies (active and inactive)
// @Tags Currency
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]currencydto.CurrencyResponse} "Successfully retrieved currencies"
// @Router /currencies [get]
func (ch *CurrencyHandler) GetAllCurrencies(c *gin.Context) {
	ctx := c.Request.Context()

	currencies, err := ch.currencyUsecase.GetAllCurrencies(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting all currencies", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get currencies")
		return
	}

	var currencyResponses []currencydto.CurrencyResponse
	for _, currency := range currencies {
		currencyResponses = append(currencyResponses, currencydto.ToCurrencyResponse(&currency))
	}

	response.Success(c, currencyResponses, "Successfully retrieved currencies")
}

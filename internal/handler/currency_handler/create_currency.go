package currency_handler

import (
	"net/http"
	"wtm-backend/internal/dto/currencydto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateCurrency godoc
// @Summary Create Currency
// @Description Create a new currency
// @Tags Currency
// @Accept json
// @Produce json
// @Param request body currencydto.CreateCurrencyRequest true "Create Currency Request"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=currencydto.CurrencyResponse} "Successfully created currency"
// @Router /currencies [post]
func (ch *CurrencyHandler) CreateCurrency(c *gin.Context) {
	ctx := c.Request.Context()

	var req currencydto.CreateCurrencyRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	currency := currencydto.ToCurrencyEntity(&req)
	created, err := ch.currencyUsecase.CreateCurrency(ctx, currency)
	if err != nil {
		logger.Error(ctx, "Error creating currency", err.Error())
		if utils.ParseValidationErrors(err) != nil {
			response.ValidationError(c, utils.ParseValidationErrors(err))
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, currencydto.ToCurrencyResponse(created), "Successfully created currency")
}

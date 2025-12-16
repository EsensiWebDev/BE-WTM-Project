package currencydto

import "wtm-backend/internal/domain/entity"

type ListCurrenciesResponse struct {
	Currencies []CurrencyResponse `json:"currencies"`
}

type CurrencyResponse struct {
	ID         uint   `json:"id"`
	ExternalID string `json:"external_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	IsActive   bool   `json:"is_active"`
}

type CreateCurrencyRequest struct {
	Code     string `json:"code" binding:"required,min=3,max=3"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Symbol   string `json:"symbol" binding:"max=10"`
	IsActive bool   `json:"is_active"`
}

type UpdateCurrencyRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Symbol   string `json:"symbol" binding:"max=10"`
	IsActive bool   `json:"is_active"`
}

func ToCurrencyResponse(currency *entity.Currency) CurrencyResponse {
	return CurrencyResponse{
		ID:         currency.ID,
		ExternalID: currency.ExternalID,
		Code:       currency.Code,
		Name:       currency.Name,
		Symbol:     currency.Symbol,
		IsActive:   currency.IsActive,
	}
}

func ToCurrencyEntity(req *CreateCurrencyRequest) *entity.Currency {
	return &entity.Currency{
		Code:     req.Code,
		Name:     req.Name,
		Symbol:   req.Symbol,
		IsActive: req.IsActive,
	}
}

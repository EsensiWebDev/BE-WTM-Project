package bookingdto

import validation "github.com/go-ozzo/ozzo-validation"

// AddGuestsToCartRequest represents the request payload for adding guests to a cart.
type AddGuestsToCartRequest struct {
	Guests []string `json:"guests"`
	CartID uint     `json:"cart_id"`
}

func (r *AddGuestsToCartRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Guests, validation.Required.Error("Guests are required")),
		validation.Field(&r.CartID, validation.Required.Error("Cart ID is required")),
	)
}

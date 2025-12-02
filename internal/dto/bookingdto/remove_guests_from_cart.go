package bookingdto

import validation "github.com/go-ozzo/ozzo-validation"

// RemoveGuestsFromCartRequest represents a request to remove guests from a cart.
type RemoveGuestsFromCartRequest struct {
	Guest  []string `json:"guest"`
	CartID uint     `json:"cart_id"`
}

func (r *RemoveGuestsFromCartRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Guest, validation.Required.Error("Guest is required")),
		validation.Field(&r.CartID, validation.Required.Error("Cart ID is required")),
	)
}

package bookingdto

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

// RemoveGuestsFromCartRequest represents a request to remove guests from a cart.
// Uses composite key (name, honorific, category, age) to uniquely identify guests.
type RemoveGuestsFromCartRequest struct {
	Guests []GuestInfo `json:"guests"` // Array of guest info to remove (composite key)
	CartID uint        `json:"cart_id"`
}

func (r *RemoveGuestsFromCartRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Guests, validation.Required.Error("Guests are required"), validation.Length(1, 1000).Error("At least one guest is required")),
		validation.Field(&r.CartID, validation.Required.Error("Cart ID is required")),
	); err != nil {
		return err
	}

	// Validate each guest
	for i, guest := range r.Guests {
		if err := guest.Validate(); err != nil {
			return fmt.Errorf("guest at index %d: %w", i, err)
		}
	}

	return nil
}

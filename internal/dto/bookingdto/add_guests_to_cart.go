package bookingdto

import (
	"fmt"
	"wtm-backend/pkg/constant"

	validation "github.com/go-ozzo/ozzo-validation"
)

// GuestInfo represents a guest with their details
type GuestInfo struct {
	Name      string `json:"name"`
	Honorific string `json:"honorific"`
	Category  string `json:"category"`      // "Adult" or "Child"
	Age       *int   `json:"age,omitempty"` // nullable, required when category="Child"
}

func (g *GuestInfo) Validate() error {
	var errs validation.Errors = make(map[string]error)

	// Validate Name
	if err := validation.ValidateStruct(g,
		validation.Field(&g.Name, validation.Required.Error("Name is required")),
		validation.Field(&g.Honorific, validation.Required.Error("Honorific is required")),
		validation.Field(&g.Category, validation.Required.Error("Category is required"),
			validation.In(constant.GuestCategoryAdult, constant.GuestCategoryChild).Error("Category must be either 'Adult' or 'Child'")),
	); err != nil {
		return err
	}

	// Validate Age based on category
	if g.Category == constant.GuestCategoryChild {
		if g.Age == nil {
			errs["age"] = validation.NewInternalError(fmt.Errorf("age is required when category is 'Child'"))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// AddGuestsToCartRequest represents the request payload for adding guests to a cart.
type AddGuestsToCartRequest struct {
	Guests []GuestInfo `json:"guests"`
	CartID uint        `json:"cart_id"`
}

func (r *AddGuestsToCartRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Guests, validation.Required.Error("Guests are required"), validation.Length(1, 0).Error("At least one guest is required")),
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

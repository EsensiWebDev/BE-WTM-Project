package bookingdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// UpdateCartAdditionalNotesRequest represents a request to update
// additional_notes for a single cart detail (booking_detail) item.
type UpdateCartAdditionalNotesRequest struct {
	SubCartID       uint   `json:"sub_cart_id"`
	AdditionalNotes string `json:"additional_notes"`
}

func (r *UpdateCartAdditionalNotesRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.SubCartID, validation.Required.Error("Sub Cart Id is required")),
		// Optional but must not exceed 500 characters (frontend already trims)
		validation.Field(&r.AdditionalNotes, validation.RuneLength(0, 500).Error("Additional Notes must not exceed 500 characters")),
	)
}



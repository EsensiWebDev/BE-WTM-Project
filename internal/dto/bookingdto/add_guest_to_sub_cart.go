package bookingdto

import validation "github.com/go-ozzo/ozzo-validation"

type AddGuestToSubCartRequest struct {
	Guest     string `json:"guest"`
	SubCartID uint   `json:"sub_cart_id"`
}

func (r *AddGuestToSubCartRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Guest, validation.Required.Error("Guest is required")),
		validation.Field(&r.SubCartID, validation.Required.Error("Sub Cart ID is required")),
	)
}

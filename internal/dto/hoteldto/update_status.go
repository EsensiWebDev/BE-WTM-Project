package hoteldto

import validation "github.com/go-ozzo/ozzo-validation"

type UpdateStatusRequest struct {
	HotelID uint `json:"hotel_id" form:"hotel_id"`
	Status  bool `json:"status" form:"status"`
}

func (usr *UpdateStatusRequest) Validate() error {
	return validation.ValidateStruct(usr,
		validation.Field(&usr.HotelID, validation.Required.Error("Hotel ID is required")),
	)
}

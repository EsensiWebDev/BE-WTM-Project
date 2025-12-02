package bookingdto

import validation "github.com/go-ozzo/ozzo-validation"

type CancelBookingRequest struct {
	SubBookingID string `json:"sub_booking_id" uri:"sub_booking_id"`
}

func (r *CancelBookingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.SubBookingID, validation.Required.Error("Sub Booking ID is required")))
}

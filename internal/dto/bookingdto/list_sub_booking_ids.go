package bookingdto

import validation "github.com/go-ozzo/ozzo-validation"

type ListSubBookingIDsRequest struct {
	BookingID string `uri:"booking_id" form:"booking_id"`
}

func (r *ListSubBookingIDsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.BookingID, validation.Required.Error("Booking ID is required")),
	)
}

type ListSubBookingIDsResponse struct {
	SubBookingIDs []string `json:"sub_booking_ids"`
}

package emaildto

import validation "github.com/go-ozzo/ozzo-validation"

// SendContactUsEmailRequest represents the request payload for sending a contact us email
type SendContactUsEmailRequest struct {
	Name           string `json:"name" form:"name"`
	Email          string `json:"email" form:"email"`
	Subject        string `json:"subject" form:"subject"`
	Type           string `json:"type" form:"type"`
	BookingCode    string `json:"booking_code" form:"booking_code"`
	SubBookingCode string `json:"sub_booking_code" form:"sub_booking_code"`
	Message        string `json:"message" form:"message"`
}

func (r *SendContactUsEmailRequest) Validate() error {
	// Validation logic can be added here if needed
	return validation.ValidateStruct(r,
		validation.Field(&r.Type,
			validation.Required.Error("Type is required"),
			validation.In("general", "booking").Error("Type must be one of: general, booking")),
	)
}

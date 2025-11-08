package emaildto

type SendContactUsEmailRequest struct {
	Name           string `json:"name" form:"name"`
	Email          string `json:"email" form:"email"`
	Subject        string `json:"subject" form:"subject"`
	Department     string `json:"department" form:"department"`
	Type           string `json:"type" form:"type"`
	BookingCode    string `json:"booking_code" form:"booking_code"`
	SubBookingCode string `json:"sub_booking_code" form:"sub_booking_code"`
	Message        string `json:"message" form:"message"`
}

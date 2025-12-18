package bookingdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// UpdateAdminNotesRequest represents a request to update
// admin_notes for a single booking detail item.
type UpdateAdminNotesRequest struct {
	SubBookingID string `json:"sub_booking_id"`
	AdminNotes   string `json:"admin_notes"`
}

func (r *UpdateAdminNotesRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.SubBookingID, validation.Required.Error("Sub Booking ID is required")),
		// Optional but must not exceed 500 characters
		validation.Field(&r.AdminNotes, validation.RuneLength(0, 500).Error("Admin Notes must not exceed 500 characters")),
	)
}

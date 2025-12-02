package bookingdto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpdateStatusRequest struct {
	BookingID    string `json:"booking_id"`
	SubBookingID string `json:"sub_booking_id"`
	StatusID     uint   `json:"status_id"`
	Reason       string `json:"reason"`
}

func (ur UpdateStatusRequest) Validate() error {
	var errs validation.Errors = make(map[string]error)
	if ur.BookingID == "" && ur.SubBookingID == "" {
		errs["booking_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or sub_booking_id must be provided"))
		errs["sub_booking_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or sub_booking_id must be provided"))
	}
	if ur.StatusID == 0 {
		errs["status_id"] = validation.NewInternalError(fmt.Errorf("status_id is required"))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

package bookingdto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpdateStatusBookingRequest struct {
	BookingID       uint   `json:"booking_id"`
	BookingDetailID uint   `json:"booking_detail_id"`
	StatusID        uint   `json:"status_id"`
	Reason          string `json:"reason"`
}

func (ur UpdateStatusBookingRequest) Validate() error {
	var errs validation.Errors = make(map[string]error)
	if ur.BookingID == 0 && ur.BookingDetailID == 0 {
		errs["booking_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or booking_detail_id must be provided"))
		errs["booking_detail_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or booking_detail_id must be provided"))
	}
	if ur.StatusID == 0 {
		errs["status_id"] = validation.NewInternalError(fmt.Errorf("status_id is required"))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

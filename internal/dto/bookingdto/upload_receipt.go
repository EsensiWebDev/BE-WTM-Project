package bookingdto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"mime/multipart"
	"strings"
)

type UploadReceiptRequest struct {
	BookingID       string                `json:"booking_id" form:"booking_id"`
	BookingDetailID string                `json:"booking_detail_id" form:"booking_detail_id"`
	Receipt         *multipart.FileHeader `json:"receipt" form:"receipt"`
}

func (r *UploadReceiptRequest) Validate() error {
	var errs validation.Errors = make(map[string]error)
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Receipt, validation.Required.Error("File receipt is required")),
	); err != nil {
		return err
	}

	if strings.TrimSpace(r.BookingID) == "" && strings.TrimSpace(r.BookingDetailID) == "" {
		errs["booking_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or booking_detail_id must be provided"))
		errs["booking_detail_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or booking_detail_id must be provided"))
		return errs
	}

	return nil

}

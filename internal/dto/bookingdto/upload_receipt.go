package bookingdto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"mime/multipart"
)

type UploadReceiptRequest struct {
	BookingID       uint                  `json:"booking_id"`
	BookingDetailID uint                  `json:"booking_detail_id"`
	FileReceipt     *multipart.FileHeader `json:"file_receipt"`
}

func (r *UploadReceiptRequest) Validate() error {
	var errs validation.Errors = make(map[string]error)
	if err := validation.ValidateStruct(r,
		validation.Field(&r.FileReceipt, validation.Required.Error("File receipt is required")),
	); err != nil {
		return err
	}

	if r.BookingID == 0 && r.BookingDetailID == 0 {
		errs["booking_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or booking_detail_id must be provided"))
		errs["booking_detail_id"] = validation.NewInternalError(fmt.Errorf("either booking_id or booking_detail_id must be provided"))
	}

	return errs

}

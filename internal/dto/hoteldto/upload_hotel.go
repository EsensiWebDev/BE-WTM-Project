package hoteldto

import (
	"fmt"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation"
)

type UploadHotelRequest struct {
	File *multipart.FileHeader `json:"file" form:"file"`
}

func (r *UploadHotelRequest) Validate() error {
	var errs validation.Errors = make(map[string]error)

	if r.File == nil || r.File.Size == 0 {
		errs["file"] = validation.NewInternalError(fmt.Errorf("file is required"))
	}

	if errs != nil && len(errs) > 0 {
		return errs
	}

	return nil
}

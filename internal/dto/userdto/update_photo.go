package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"mime/multipart"
)

type UpdateFileRequest struct {
	FileType string                `json:"file_type" form:"file_type"`
	File     *multipart.FileHeader `json:"photo" form:"photo"`
}

func (r *UpdateFileRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.File, validation.Required.Error("File is required")),
		validation.Field(&r.FileType,
			validation.Required.Error("File type is required"),
			validation.In("photo", "certificate", "name_card").
				Error("File type must be one of: photo, certificate, name_card")),
	)
}

package emaildto

import (
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation"
)

// UpdateEmailTemplateRequest represents the request to update an email template
type UpdateEmailTemplateRequest struct {
	Type           string                `json:"type" form:"type"`
	Subject        string                `json:"subject" form:"subject"`
	Body           string                `json:"body" form:"body"`
	SignatureText  string                `json:"signature_text" form:"signature_text"`
	SignatureImage *multipart.FileHeader `json:"signature_image" form:"signature_image"`
}

func (r *UpdateEmailTemplateRequest) Validate() error {

	return validation.ValidateStruct(r,
		validation.Field(&r.Body, validation.Required.Error("Template body is required")),
	)
}

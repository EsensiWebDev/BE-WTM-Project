package emaildto

import validation "github.com/go-ozzo/ozzo-validation"

type EmailTemplateRequest struct {
	Type string `json:"type" form:"type"`
}

func (e *EmailTemplateRequest) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Type, validation.In("", "confirm", "cancel")),
	)
}

type EmailTemplateResponse struct {
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Signature string `json:"signature"`
}

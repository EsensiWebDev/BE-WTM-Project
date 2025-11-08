package emaildto

type EmailTemplateResponse struct {
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Signature string `json:"signature"`
}

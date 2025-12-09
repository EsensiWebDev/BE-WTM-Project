package entity

type EmailTemplate struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Subject          string `json:"subject"`
	Body             string `json:"body"`
	IsSignatureImage bool   `json:"is_signature_image"`
	Signature        string `json:"signature"`
	ExternalID       string
}

type Notification struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	RedirectURL string `json:"redirect_url"`
	Type        string `json:"type"`
	IsRead      bool   `json:"is_read"`
	ReadAt      string `json:"read_at"`
}

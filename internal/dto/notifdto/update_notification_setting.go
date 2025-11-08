package notifdto

import validation "github.com/go-ozzo/ozzo-validation"

type UpdateNotificationSettingRequest struct {
	Channel  string
	Type     string
	IsEnable bool
}

func (r *UpdateNotificationSettingRequest) Validate() error {
	// Add validation logic if needed
	if r.IsEnable {
		switch r.Type {
		case "booking", "reject", "all":
			// valid type
		default:
			return validation.Validate(r.Type, validation.In("booking", "reject", "all"))
		}
	}
	return validation.ValidateStruct(r,
		validation.Field(&r.Channel, validation.Required, validation.In("email", "web")),
	)
}

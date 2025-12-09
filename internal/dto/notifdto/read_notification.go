package notifdto

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ReadNotificationRequest struct {
	Type string `json:"type" form:"type"`
	ID   uint   `json:"notif_id" form:"notif_id"`
}

func (r *ReadNotificationRequest) Validate() error {
	if r.Type != "all" {
		if r.ID == 0 {
			return validation.NewInternalError(fmt.Errorf("ID is required when type is not 'all'"))
		}
	}
	return nil
}

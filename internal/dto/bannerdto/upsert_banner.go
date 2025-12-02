package bannerdto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"mime/multipart"
	"wtm-backend/pkg/utils"
)

type UpsertBannerRequest struct {
	Title       string                `json:"title" form:"title"`
	Description string                `json:"description" form:"description"`
	Image       *multipart.FileHeader `json:"image" form:"image"`
}

func (r *UpsertBannerRequest) ValidateCreate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Title, validation.Required.Error("Title is required"), utils.NotEmptyAfterTrim("Title")),
		validation.Field(&r.Image, validation.Required.Error("Image is required")),
	); err != nil {
		return err
	}

	if r.Image == nil || r.Image.Size == 0 {
		return validation.Errors{
			"image": validation.NewInternalError(fmt.Errorf("image is required")),
		}
	}

	return nil
}

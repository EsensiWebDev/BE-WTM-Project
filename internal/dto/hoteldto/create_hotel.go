package hoteldto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"mime/multipart"
	"wtm-backend/pkg/utils"
)

type CreateHotelRequest struct {
	Name         string                  `json:"name" form:"name"`
	Photos       []*multipart.FileHeader `json:"photos" form:"photos"`
	SubDistrict  string                  `json:"sub_district" form:"sub_district"`
	District     string                  `json:"district" form:"district"`
	Email        string                  `json:"email" form:"email"`
	Province     string                  `json:"province" form:"province"`
	Description  string                  `json:"description" form:"description"`
	Rating       int                     `json:"rating" form:"rating"`
	NearbyPlaces string                  `json:"nearby_places" form:"nearby_places"`
	Facilities   []string                `json:"facilities" form:"facilities"`
	SocialMedias string                  `json:"social_medias" form:"social_medias"`
}

// CreateHotelResponse represents the response create hotels.
type CreateHotelResponse struct {
	HotelID uint `json:"hotel_id"`
}

// NearbyPlace represents the request structure nearby places.
type NearbyPlace struct {
	Name     string  `json:"name" form:"name"`
	Distance float64 `json:"distance" form:"distance"`
}

// SocialMedia represents the request structure social medias.
type SocialMedia struct {
	Platform string `json:"platform" form:"platform"`
	Link     string `json:"link" form:"link"`
}

func (r *CreateHotelRequest) Validate() error {
	var errs validation.Errors = make(map[string]error)
	if err := validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required.Error("Hotel name is required"), utils.NotEmptyAfterTrim("Name")),
		validation.Field(&r.SubDistrict, validation.Required.Error("Sub-district is required"), utils.NotEmptyAfterTrim("Sub-Cities")),
		validation.Field(&r.District, validation.Required.Error("Cities is required"), utils.NotEmptyAfterTrim("Cities")),
		validation.Field(&r.Province, validation.Required.Error("Province is required"), utils.NotEmptyAfterTrim("Province")),
	); err != nil {
		return err
	}

	if len(r.Photos) == 0 {
		errs["photos"] = validation.NewInternalError(fmt.Errorf("at least one photo is required"))
	}

	if errs != nil && len(errs) > 0 {
		return errs
	}

	return nil
}

package promodto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpsertPromoRequest struct {
	StartDate   string     `json:"start_date" form:"start_date"`
	EndDate     string     `json:"end_date" form:"end_date"`
	RoomTypes   []RoomType `json:"room_types" form:"room_types"`
	PromoName   string     `json:"promo_name" form:"promo_name"`
	PromoTypeID uint       `json:"promo_type" form:"promo_type"`
	Detail      string     `json:"detail" form:"detail"`
	PromoCode   string     `json:"promo_code" form:"promo_code"`
	Description string     `json:"description" form:"description"`
}

type RoomType struct {
	RoomTypeID uint `json:"room_type_id" form:"room_type_id"`
	TotalNight int  `json:"total_night" form:"total_night"`
}

func (r *UpsertPromoRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.StartDate,
			validation.Required.
				Error("Start date is required"),
			validation.Date("2006-01-02T15:04:05Z07:00").
				Error("Start date must be in RFC3339 format (e.g., 2023-10-01T00:00:00Z)"),
		),
		validation.Field(&r.EndDate,
			validation.Required.
				Error("End date is required"),
			validation.Date("2006-01-02T15:04:05Z07:00").
				Error("End date must be in RFC3339 format (e.g., 2023-10-01T00:00:00Z)"),
		),
		validation.Field(&r.PromoName, validation.Required.Error("Promo name is required")),
		validation.Field(&r.PromoTypeID, validation.Required.Error("Promo type Id is required")),
		validation.Field(&r.Detail, validation.Required.Error("Detail is required")),
		validation.Field(&r.PromoCode, validation.Required.Error("Promo code is required")),
	); err != nil {
		return err
	}

	if len(r.RoomTypes) == 0 {
		return validation.Errors{
			"room_types": validation.NewInternalError(fmt.Errorf("at least one room type is required")),
		}
	}

	return nil
}

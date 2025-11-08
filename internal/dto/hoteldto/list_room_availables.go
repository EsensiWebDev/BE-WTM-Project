package hoteldto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/pkg/utils"
)

type ListRoomAvailableRequest struct {
	HotelID uint   `json:"hotel_id" form:"hotel_id"`
	Month   string `json:"month" form:"month"`
}

type ListRoomAvailableResponse struct {
	RoomAvailable []RoomAvailable `json:"room_unavailable"`
}

type RoomAvailable struct {
	RoomTypeID   uint            `json:"room_type_id"`
	RoomTypeName string          `json:"room_type_name"`
	Data         []DataAvailable `json:"available"`
}

type DataAvailable struct {
	Day       int  `json:"day"`
	Available bool `json:"available"`
}

func (r *ListRoomAvailableRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.HotelID, validation.Required.Error("Hotel Id is required")),
		validation.Field(&r.Month, validation.Required.Error("Month is required")),
	); err != nil {
		return err
	}

	if !utils.IsValidMonth(r.Month) {
		return validation.Errors{
			"month": validation.NewInternalError(fmt.Errorf("invalid month format, must be YYYY-MM")),
		}
	}

	return nil
}

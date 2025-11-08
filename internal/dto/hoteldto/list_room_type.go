package hoteldto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/dto"
)

type ListRoomTypeRequest struct {
	HotelID               uint `json:"hotel_id" form:"hotel_id"`
	dto.PaginationRequest `json:",inline"`
}

type ListRoomTypeResponse struct {
	RoomTypes []ListRoomType `json:"room_types"`
}

type ListRoomType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (r *ListRoomTypeRequest) Validate() error {
	return validation.ValidateStruct(r, validation.Field(&r.HotelID, validation.Required.Error("Hotel Id is required")))
}

package hoteldto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/domain/entity"
)

type ListBedTypeRequest struct {
	RoomTypeID uint `json:"room_type_id" form:"room_type_id"`
}

type ListBedTypeResponse struct {
	BedTypes []entity.BedType `json:"bed_types"`
}

func (r *ListBedTypeRequest) Validate() error {
	return validation.ValidateStruct(r, validation.Field(&r.RoomTypeID, validation.Required.Error("RoomTypeID is required")))
}

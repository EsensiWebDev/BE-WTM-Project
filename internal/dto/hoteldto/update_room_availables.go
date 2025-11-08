package hoteldto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/pkg/utils"
)

type UpdateRoomAvailableRequest struct {
	Month string                    `json:"month"` // Format: YYYY-MM
	Data  []UpdateRoomAvailableData `json:"data"`
}

type UpdateRoomAvailableData struct {
	RoomTypeID    uint            `json:"room_type_id"`
	RoomAvailable []DataAvailable `json:"room_available"`
}

func (r *UpdateRoomAvailableRequest) Validate() error {
	if !utils.IsValidMonth(r.Month) {
		return validation.Errors{
			"month": validation.NewInternalError(fmt.Errorf("invalid month format, must be YYYY-MM")),
		}
	}
	if len(r.Data) == 0 {
		return validation.Errors{
			"data": validation.NewInternalError(fmt.Errorf("at least one room type data is required")),
		}
	}

	for _, data := range r.Data {
		if data.RoomTypeID == 0 {
			return validation.Errors{
				"room_type_id": validation.NewInternalError(fmt.Errorf("room type Id is required")),
			}
		}

		if len(data.RoomAvailable) == 0 {
			return validation.Errors{
				"room_available": validation.NewInternalError(fmt.Errorf("at least one available date is required for room type %d", data.RoomTypeID)),
			}
		}
	}

	return nil
}

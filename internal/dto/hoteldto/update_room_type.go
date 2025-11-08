package hoteldto

import "mime/multipart"

type UpdateRoomTypeRequest struct {
	RoomTypeID            uint                    `json:"room_type_id" form:"room_type_id"`
	Name                  string                  `json:"name" form:"name"`
	Photos                []*multipart.FileHeader `json:"photos" form:"photos"`
	WithoutBreakfast      string                  `json:"without_breakfast" form:"without_breakfast"`
	WithBreakfast         string                  `json:"with_breakfast" form:"with_breakfast"`
	RoomSize              float64                 `json:"room_size" form:"room_size"`
	MaxOccupancy          int                     `json:"max_occupancy" form:"max_occupancy"`
	BedTypes              []string                `json:"bed_types" form:"bed_types"`
	IsSmokingRoom         bool                    `json:"is_smoking_room" form:"is_smoking_room"`
	Additional            string                  `json:"additional" form:"additional"`
	Description           string                  `json:"description" form:"description"`
	UnchangedRoomPhotos   []string                `json:"unchanged_room_photos" form:"unchanged_room_photos"`
	UnchangedAdditionsIDs []uint                  `json:"unchanged_additions_ids" form:"unchanged_additions_ids"`
}

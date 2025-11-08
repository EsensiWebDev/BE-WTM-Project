package hoteldto

import (
	"mime/multipart"
)

// AddRoomTypeRequest represents the request structure for adding a new room type to a hotel.
type AddRoomTypeRequest struct {
	HotelID          uint                    `json:"hotel_id" form:"hotel_id"`
	Name             string                  `json:"name" form:"name"`
	Photos           []*multipart.FileHeader `json:"photos" form:"photos"`
	WithoutBreakfast string                  `json:"without_breakfast" form:"without_breakfast"`
	WithBreakfast    string                  `json:"with_breakfast" form:"with_breakfast"`
	RoomSize         float64                 `json:"room_size" form:"room_size"`
	MaxOccupancy     int                     `json:"max_occupancy" form:"max_occupancy"`
	BedTypes         []string                `json:"bed_types" form:"bed_types"`
	IsSmokingRoom    bool                    `json:"is_smoking_room" form:"is_smoking_room"`
	Additional       string                  `json:"additional" form:"additional"`
	Description      string                  `json:"description" form:"description"`
	//TotalUnit        int                     `json:"total_unit" form:"total_unit"`
}

type BreakfastBase struct {
	Price  float64 `json:"price" form:"price"`
	IsShow bool    `json:"is_show" form:"is_show"`
}

type BreakfastWith struct {
	BreakfastBase
	Pax int `json:"pax" form:"pax"`
}

type RoomAdditional struct {
	Name  string  `json:"name" form:"name"`
	Price float64 `json:"price" form:"price"`
}

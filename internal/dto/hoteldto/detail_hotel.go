package hoteldto

import (
	"wtm-backend/internal/domain/entity"
)

// DetailHotelResponse represents the detailed information of a hotel.
type DetailHotelResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Province    string               `json:"province"`
	District    string               `json:"city"`
	SubDistrict string               `json:"sub_district"`
	Description string               `json:"description"`
	Photos      []string             `json:"photos"`
	Rating      int                  `json:"rating"`
	Email       string               `json:"email"`
	Facilities  []string             `json:"facilities"`
	NearbyPlace []entity.NearbyPlace `json:"nearby_place"`
	SocialMedia []SocialMedia        `json:"social_media"`
	RoomType    []DetailRoomType     `json:"room_type"`

	CancellationPeriod int    `json:"cancellation_period"`
	CheckInHour        string `json:"check_in_hour"`
	CheckOutHour       string `json:"check_out_hour"`
}

type DetailRoomType struct {
	ID               uint                                `json:"id"`
	Name             string                              `json:"name"`
	WithoutBreakfast entity.CustomBreakfast              `json:"without_breakfast"`
	WithBreakfast    entity.CustomBreakfast              `json:"with_breakfast"`
	RoomSize         float64                             `json:"room_size"`
	MaxOccupancy     int                                 `json:"max_occupancy"`
	BedTypes         []string                            `json:"bed_types"`
	IsSmokingRoom    bool                                `json:"is_smoking_room"`
	Additional       []entity.CustomRoomAdditionalWithID `json:"additional"`
	Description      string                              `json:"description"`
	Photos           []string                            `json:"photos"`
	//TotalUnit        int                           `json:"total_unit"`
}

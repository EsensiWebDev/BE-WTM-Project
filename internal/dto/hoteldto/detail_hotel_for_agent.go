package hoteldto

import "wtm-backend/internal/domain/entity"

type DetailHotelForAgentResponse struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"name"`
	Province    string                   `json:"province"`
	District    string                   `json:"city"`
	SubDistrict string                   `json:"sub_district"`
	Description string                   `json:"description"`
	Photos      []string                 `json:"photos"`
	Rating      int                      `json:"rating"`
	Email       string                   `json:"email"`
	Facilities  []string                 `json:"facilities"`
	NearbyPlace []entity.NearbyPlace     `json:"nearby_place"`
	SocialMedia []SocialMedia            `json:"social_media"`
	RoomType    []DetailRoomTypeForAgent `json:"room_type"`

	CancellationPeriod int    `json:"cancellation_period"`
	CheckInHour        string `json:"check_in_hour"`
	CheckOutHour       string `json:"check_out_hour"`
}

type DetailRoomTypeForAgent struct {
	Name             string                              `json:"name"`
	WithoutBreakfast entity.CustomBreakfastWithID        `json:"without_breakfast"`
	WithBreakfast    entity.CustomBreakfastWithID        `json:"with_breakfast"`
	RoomSize         float64                             `json:"room_size"`
	MaxOccupancy     int                                 `json:"max_occupancy"`
	BedTypes         []string                            `json:"bed_types"`
	IsSmokingRoom    bool                                `json:"is_smoking_room"`
	Additional       []entity.CustomRoomAdditionalWithID `json:"additional"`
	Description      string                              `json:"description"`
	Photos           []string                            `json:"photos"`
	Promos           []PromoDetailRoom                   `json:"promos"`
}

type PromoDetailRoom struct {
	PromoID               uint    `json:"promo_id"`
	Description           string  `json:"description"`
	CodePromo             string  `json:"code_promo"`
	PriceWithBreakfast    float64 `json:"price_with_breakfast"`
	PriceWithoutBreakfast float64 `json:"price_without_breakfast"`
	TotalNights           int     `json:"total_nights"`
	OtherNotes            string  `json:"other_notes"`
}

package hoteldto

type UpdateHotelRequest struct {
	CreateHotelRequest      `json:",inline" form:",inline"`
	HotelID                 uint     `json:"hotel_id" form:"hotel_id"`
	UnchangedHotelPhotos    []string `json:"unchanged_hotel_photos" form:"unchanged_hotel_photos"`
	UnchangedNearbyPlaceIDs []uint   `json:"unchanged_nearby_place_ids" form:"unchanged_nearby_place_ids"`
}

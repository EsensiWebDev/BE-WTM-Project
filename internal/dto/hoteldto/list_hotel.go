package hoteldto

import "wtm-backend/internal/dto"

type ListHotelRequest struct {
	dto.PaginationRequest `json:",inline"`
	IsAPI                 *bool    `json:"is_api,omitempty" form:"is_api,omitempty"`
	Region                []string `json:"region" form:"region"`
	StatusID              uint     `json:"status_id" form:"status_id"`
}

type ListHotelResponse struct {
	Hotels []ListHotel `json:"hotels"`
	Total  int64       `json:"total"`
}

type ListHotel struct {
	ID     uint           `json:"id"`
	Name   string         `json:"name"`
	Region string         `json:"region"`
	Email  string         `json:"email"`
	Status string         `json:"status"`
	IsAPI  bool           `json:"is_api"`
	Rooms  []RoomTypeItem `json:"rooms"`
}

type RoomTypeItem struct {
	Name                  string  `json:"name"`
	Price                 float64 `json:"price"`
	PriceWithoutBreakfast float64 `json:"price_with_breakfast"`
}

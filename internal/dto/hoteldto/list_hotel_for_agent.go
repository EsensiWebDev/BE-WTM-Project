package hoteldto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListHotelForAgentRequest struct {
	dto.PaginationRequest `json:",inline"`
	Province              *string  `json:"province" form:"province"`
	Rating                []int    `json:"rating" form:"rating"`
	BedTypeID             []int    `json:"bed_type_id" form:"bed_type_id"`
	RangePriceMin         *int     `json:"range_price_min" form:"range_price_min"`
	RangePriceMax         *int     `json:"range_price_max" form:"range_price_max"`
	District              []string `json:"district" form:"district"`
	TotalBedrooms         []int    `json:"total_bedrooms" form:"total_bedrooms"`
	TotalRooms            int      `json:"total_rooms" form:"total_rooms"`
	RangeDateFrom         string   `json:"from" form:"from"`
	RangeDateTo           string   `json:"to" form:"to"`
	TotalGuests           int      `json:"total_guests" form:"total_guests"`
	PromoID               int64    `json:"promo_id" form:"promo_id"`
}

type ListHotelForAgentResponse struct {
	Hotels           []ListHotelForAgent         `json:"hotels"`
	FilterDistricts  []string                    `json:"filter_districts"`
	FilterPricing    *entity.FilterRangePrice    `json:"filter_pricing"`
	FilterRatings    []entity.FilterRatingHotel  `json:"filter_ratings"`
	FilterBedTypes   []entity.FilterBedTypeHotel `json:"filter_bed_types"`
	FilterTotalRooms []entity.FilterTotalBedroom `json:"filter_total_rooms"`
	Total            int64                       `json:"total"`
}

type ListHotelForAgent struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	MinPrice float64 `json:"min_price"`
	Photo    string  `json:"photo"`
	Rating   int     `json:"rating"`
}

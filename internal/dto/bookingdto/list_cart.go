package bookingdto

import "time"

type ListCartResponse struct {
	ID         uint         `json:"id"`
	Detail     []CartDetail `json:"detail"`
	Guest      []string     `json:"guest"`
	GrandTotal int64        `json:"grand_total"`
}

type CartDetail struct {
	HotelName            string                 `json:"hotel_name"`
	HotelRating          int                    `json:"hotel_rating"`
	CheckInDate          time.Time              `json:"check_in_date"`
	CheckOutDate         time.Time              `json:"check_out_date"`
	RoomTypeName         string                 `json:"room_type_name"`
	IsBreakfast          bool                   `json:"is_breakfast"`
	Guest                string                 `json:"guest"`
	Additional           []CartDetailAdditional `json:"additional"`
	Promo                CartDetailPromo        `json:"promo"`
	Price                float64                `json:"price"`
	TotalAdditionalPrice float64                `json:"total_additional_price"`
	TotalPrice           float64                `json:"total_price"`
}

type CartDetailAdditional struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CartDetailPromo struct {
	Type            string  `json:"type,omitempty"`
	Code            string  `json:"code,omitempty"`
	DiscountPercent float64 `json:"discount_percent,omitempty"`
	FixedPrice      float64 `json:"fixed_price,omitempty"`
	UpgradedToID    uint    `json:"upgraded_to_id,omitempty"`
	Benefit         string  `json:"benefit,omitempty"`
}

package bookingdto

import (
	"time"
	"wtm-backend/internal/domain/entity"
)

type ListCartResponse struct {
	ID         uint         `json:"id"`
	Detail     []CartDetail `json:"detail"`
	Guest      []string     `json:"guest"`
	GrandTotal float64      `json:"grand_total"`
}

type CartDetail struct {
	ID                   uint                   `json:"id"`
	Photo                string                 `json:"photo"`
	HotelName            string                 `json:"hotel_name"`
	HotelRating          int                    `json:"hotel_rating"`
	CheckInDate          time.Time              `json:"check_in_date"`
	CheckOutDate         time.Time              `json:"check_out_date"`
	RoomTypeName         string                 `json:"room_type_name"`
	IsBreakfast          bool                   `json:"is_breakfast"`
	Guest                string                 `json:"guest"`
	Additional           []CartDetailAdditional `json:"additional"`
	Promo                entity.DetailPromo     `json:"promo"`
	CancellationDate     string                 `json:"cancellation_date,omitempty"`
	Price                float64                `json:"price"`
	PriceBeforePromo     float64                `json:"price_before_promo"`
	TotalAdditionalPrice float64                `json:"total_additional_price"`
	TotalPrice           float64                `json:"total_price"`
}

type CartDetailAdditional struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

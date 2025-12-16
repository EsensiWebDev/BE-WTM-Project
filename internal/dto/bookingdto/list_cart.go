package bookingdto

import (
	"wtm-backend/internal/domain/entity"
)

type ListCartResponse struct {
	ID         uint         `json:"id"`
	Detail     []CartDetail `json:"detail"`
	Guest      []CartGuest  `json:"guest"`
	GrandTotal float64      `json:"grand_total"`
}

type CartGuest struct {
	Name      string `json:"name"`
	Honorific string `json:"honorific"`
	Category  string `json:"category"`      // "Adult" or "Child"
	Age       *int   `json:"age,omitempty"` // nullable, only present when category="Child"
}

type CartDetail struct {
	ID                   uint                   `json:"id"`
	Photo                string                 `json:"photo"`
	HotelName            string                 `json:"hotel_name"`
	HotelRating          int                    `json:"hotel_rating"`
	CheckInDate          string                 `json:"check_in_date"`
	CheckOutDate         string                 `json:"check_out_date"`
	RoomTypeName         string                 `json:"room_type_name"`
	IsBreakfast          bool                   `json:"is_breakfast"`
	Guest                string                 `json:"guest"`
	BedType              string                 `json:"bed_type,omitempty"`  // Selected bed type (singular) - REQUIRED
	BedTypes             []string               `json:"bed_types,omitempty"` // Available bed types for reference (plural) - OPTIONAL
	OtherPreferences     []string               `json:"other_preferences"`
	Additional           []CartDetailAdditional `json:"additional"`
	Promo                entity.DetailPromo     `json:"promo"`
	CancellationDate     string                 `json:"cancellation_date,omitempty"`
	Price                float64                `json:"price"`
	PriceBeforePromo     float64                `json:"price_before_promo"`
	TotalAdditionalPrice float64                `json:"total_additional_price"`
	TotalPrice           float64                `json:"total_price"`
}

type CartDetailAdditional struct {
	Name       string   `json:"name"`
	Category   string   `json:"category"`        // "price" or "pax"
	Price      *float64 `json:"price,omitempty"` // nullable, used when category="price"
	Pax        *int     `json:"pax,omitempty"`   // nullable, used when category="pax"
	IsRequired bool     `json:"is_required"`
}

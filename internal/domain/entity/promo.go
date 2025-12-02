package entity

import "time"

// PromoGroup represents a group of promotions
type PromoGroup struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}

type Promo struct {
	ID          uint        `json:"id"`
	ExternalID  string      `json:"external_id"`
	Name        string      `json:"name"`
	StartDate   *time.Time  `json:"start_date,omitempty"`
	EndDate     *time.Time  `json:"end_date,omitempty"`
	Code        string      `json:"code,omitempty"`
	Description string      `json:"description,omitempty"`
	PromoTypeID uint        `json:"promo_type_id,omitempty"`
	Detail      PromoDetail `json:"detail,omitempty"`
	IsActive    bool        `json:"is_active,omitempty"`

	PromoTypeName string `json:"promo_type_name,omitempty"`

	PromoGroups []PromoGroup `json:"promo_groups,omitempty"`

	PromoRoomTypes []PromoRoomType `json:"promo_room_types,omitempty"`

	PromoGroupIDs []uint `json:"promo_group_ids,omitempty"`
	Duration      int    `json:"duration,omitempty"`
}

type PromoDetail struct {
	DiscountPercentage float64 `json:"discount_percentage,omitempty"`
	FixedPrice         float64 `json:"fixed_price,omitempty"`
	UpgradedToID       uint    `json:"upgraded_to_id,omitempty"`
	BenefitNote        string  `json:"benefit_note,omitempty"`
}

type PromoRoomType struct {
	RoomTypeID   uint   `json:"room_type_id"`
	RoomTypeName string `json:"room_type_name"`
	TotalNights  int    `json:"total_nights"`
	HotelID      uint   `json:"hotel_id"`
	HotelName    string `json:"hotel_name"`
	Province     string `json:"province"`
}

type PromoType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Banner struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	Order       int    `json:"order"`
	ExternalID  string `json:"external_id"`
}

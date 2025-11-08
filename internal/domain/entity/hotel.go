package entity

import (
	"time"
)

type Hotel struct {
	ID              uint
	Name            string
	AddrSubDistrict string
	AddrCity        string
	AddrProvince    string
	IsAPI           bool
	UrlAPI          string
	Description     string
	Photos          []string
	StatusID        uint
	Rating          int
	Email           string

	StatusHotel   string
	FacilityNames []string
	NearbyPlaces  []NearbyPlace
	RoomTypes     []RoomType
	PromoHotel    []Promo

	CancellationPeriod int
	CheckInHour        *time.Time
	CheckOutHour       *time.Time

	SocialMedia map[string]string
}

type CustomHotel struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	AddrSubDistrict string   `json:"addr_sub_district"`
	AddrCity        string   `json:"addr_city"`
	AddrProvince    string   `json:"addr_province"`
	Photos          []string `json:"photos"`
	Rating          int      `json:"rating"`
	MinPrice        float64  `json:"min_price"`
}

type BedType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type NearbyPlace struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Radius float64 `json:"radius"`
}

type RoomType struct {
	ID               uint
	HotelID          uint
	Name             string
	IsSmokingAllowed *bool
	MaxOccupancy     int
	RoomSize         float64
	Description      string
	Photos           []string

	WithoutBreakfast CustomBreakfastWithID
	WithBreakfast    CustomBreakfastWithID

	BedTypeNames   []string
	RoomAdditions  []CustomRoomAdditionalWithID
	Hotel          Hotel
	PromoRoomTypes []PromoRoomTypes

	TotalUnit int
}

type PromoRoomTypes struct {
	ID         uint
	PromoID    uint
	RoomTypeID uint
	TotalNight int
	Promo      Promo
}

type CustomRoomAdditional struct {
	Name  string  `json:"name" form:"name"`
	Price float64 `json:"price" form:"price"`
}

type CustomRoomAdditionalWithID struct {
	ID    uint    `json:"id" form:"id"`
	Name  string  `json:"name" form:"name"`
	Price float64 `json:"price" form:"price"`
}

type CustomBreakfast struct {
	Price  float64 `json:"price" form:"price"`
	Pax    int     `json:"pax,omitempty" form:"pax"`
	IsShow bool    `json:"is_show" form:"is_show"`
}

type CustomBreakfastWithID struct {
	ID     uint    `json:"id" form:"id"`
	Price  float64 `json:"price" form:"price"`
	Pax    int     `json:"pax,omitempty" form:"pax"`
	IsShow bool    `json:"is_show" form:"is_show"`
}

type RoomUnavailable struct {
	RoomTypeID uint
	Date       *time.Time
}

type FilterRangePrice struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type FilterRatingHotel struct {
	Rating int `json:"rating"`
	Count  int `json:"count"`
}

type FilterBedTypeHotel struct {
	BedTypeID uint   `json:"bed_type_id"`
	BedType   string `json:"bed_type"`
	Count     int    `json:"count"`
}

type FilterTotalBedroom struct {
	TotalBedRooms int `json:"total_bed_rooms"`
	Count         int `json:"count"`
}

type RoomPrice struct {
	ID          uint
	RoomTypeID  uint
	IsBreakfast bool
	Pax         int
	Price       float64
	IsShow      bool

	RoomType RoomType
}

type RoomTypeAdditional struct {
	ID               uint
	RoomTypeID       uint
	RoomAdditionalID uint
	Price            float64

	RoomAdditional RoomAdditional
}

type RoomAdditional struct {
	ID   uint
	Name string
}

type StatusHotel struct {
	ID     uint
	Status string
}

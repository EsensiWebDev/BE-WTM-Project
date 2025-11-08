package model

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type StatusHotel struct {
	ID     uint   `json:"id" gorm:"primaryKey"` // override default gorm.Model ID
	Status string `json:"status"`
}

type Hotel struct {
	gorm.Model
	Name            string         `json:"name"`
	AddrSubDistrict string         `json:"addr_sub_district"`
	AddrCity        string         `json:"addr_city"`
	AddrProvince    string         `json:"addr_province"`
	IsAPI           bool           `json:"is_api"`
	UrlAPI          string         `json:"url_api"`
	Description     string         `json:"description"`
	Photos          pq.StringArray `json:"photos" gorm:"type:text[]"`
	StatusID        uint           `json:"status_id" gorm:"index;default:1"`
	Rating          int            `json:"rating" gorm:"default:0"`
	Email           string         `json:"email" gorm:"uniqueIndex"`

	CancellationPeriod int        `json:"cancellation_period" gorm:"default:0"`
	CheckInHour        *time.Time `json:"check_in_hour" gorm:"default:null;type:time"`
	CheckOutHour       *time.Time `json:"check_out_hour" gorm:"default:null;type:time"`

	SocialMedia datatypes.JSON `json:"social_media" gorm:"type:jsonb"`

	Status     StatusHotel `gorm:"foreignkey:StatusID"`
	Facilities []Facility  `gorm:"many2many:HotelFacility"`

	HotelNearbyPlaces []HotelNearbyPlace `gorm:"foreignKey:HotelID"`
	RoomTypes         []RoomType         `gorm:"foreignKey:HotelID"`
}

type NearbyPlace struct {
	gorm.Model
	Name string `json:"name"`

	Hotel []Hotel `gorm:"many2many:HotelNearbyPlace"`
}

type HotelNearbyPlace struct {
	gorm.Model
	HotelID       uint    `json:"hotel_id" gorm:"index"`
	NearbyPlaceID uint    `json:"nearby_place_id" gorm:"index"`
	Radius        float64 `json:"radius"`

	Hotel       Hotel       `json:"hotel" gorm:"foreignkey:HotelID"`
	NearbyPlace NearbyPlace `json:"nearby_place" gorm:"foreignkey:NearbyPlaceID"`
}

type Facility struct {
	gorm.Model
	Name string `json:"name"`

	Hotels []Hotel `gorm:"many2many:HotelFacility"`
}

package filter

import (
	"strings"
	"time"
	"wtm-backend/internal/dto"
)

type PromoGroupFilter struct {
	PromoGroupID uint
	dto.PaginationRequest
}

type UserFilter struct {
	RoleID         *uint
	AgentCompanyID *uint
	Scope          string
	StatusID       *uint
	dto.PaginationRequest
}

type HotelFilter struct {
	IsAPI    *bool
	Region   string
	StatusID uint
	dto.PaginationRequest
}

type DefaultFilter struct {
	dto.PaginationRequest
}

type BannerFilter struct {
	dto.PaginationRequest
	IsActive bool
}

type ReportFilter struct {
	DateFrom       *time.Time
	DateTo         *time.Time
	HotelID        *uint
	AgentCompanyID *uint
	IsRangeDate    bool
	dto.PaginationRequest
}

type ReportDetailFilter struct {
	HotelID *uint
	AgentID *uint
	dto.PaginationRequest
}

type HotelFilterForAgent struct {
	Ratings       []int
	BedTypeIDs    []uint
	PriceMin      *int
	PriceMax      *int
	Cities        []string
	TotalBedrooms []int

	Province   *string
	DateFrom   *time.Time
	DateTo     *time.Time
	TotalGuest *int
	dto.PaginationRequest
}

type BookingFilter struct {
	dto.PaginationRequest
	AgentID         uint
	BookingIDSearch string
	GuestNameSearch string
	StatusBookingID uint
	StatusPaymentID uint
}

type NotifFilter struct {
	dto.PaginationRequest
	UserID uint
}

func Clean[T any](filter *T) *T {
	switch v := any(filter).(type) {
	case *HotelFilterForAgent:
		cleanHotelFilter(v)
	case *UserFilter:
		cleanUserFilter(v)
	}
	return filter
}

func cleanHotelFilter(f *HotelFilterForAgent) {
	f.Ratings = cleanIntSlice(f.Ratings)
	f.BedTypeIDs = cleanUintSlice(f.BedTypeIDs)
	if f.PriceMin != nil && *f.PriceMin <= 0 {
		f.PriceMin = nil
	}
	if f.PriceMax != nil && *f.PriceMax <= 0 {
		f.PriceMax = nil
	}
	f.Cities = cleanStringSlice(f.Cities)
	f.TotalBedrooms = cleanIntSlice(f.TotalBedrooms)
}

func cleanUserFilter(f *UserFilter) {
	if f.RoleID != nil && *f.RoleID == 0 {
		f.RoleID = nil
	}
	if f.AgentCompanyID != nil && *f.AgentCompanyID == 0 {
		f.AgentCompanyID = nil
	}
}

func cleanIntSlice(input []int) []int {
	if len(input) == 0 {
		return nil
	}

	var out []int
	for _, v := range input {
		if v > 0 {
			out = append(out, v)
		}
	}

	return out
}

func cleanUintSlice(input []uint) []uint {
	if len(input) == 0 {
		return nil
	}

	var out []uint
	for _, v := range input {
		if v > 0 {
			out = append(out, v)
		}
	}
	return out
}

func cleanStringSlice(input []string) []string {
	if len(input) == 0 {
		return nil
	}
	var out []string
	for _, v := range input {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out
}

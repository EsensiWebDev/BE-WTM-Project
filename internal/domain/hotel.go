package domain

import (
	"context"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
)

type HotelUsecase interface {
	ListHotels(ctx context.Context, req *hoteldto.ListHotelRequest) (*hoteldto.ListHotelResponse, error)
	ListHotelsForAgent(ctx context.Context, req *hoteldto.ListHotelForAgentRequest) (*hoteldto.ListHotelForAgentResponse, error)
	ListRoomTypes(ctx context.Context, hotelID uint) (*hoteldto.ListRoomTypeResponse, error)
	ListBedTypes(ctx context.Context, roomTypeID uint) (*hoteldto.ListBedTypeResponse, error)
	CreateHotel(ctx context.Context, req *hoteldto.CreateHotelRequest) error
	ListFacilities(ctx context.Context, req *hoteldto.ListFacilitiesRequest) (*hoteldto.ListFacilitiesResponse, error)
	ListAdditionalRooms(ctx context.Context, req *hoteldto.ListAdditionalRoomsRequest) (*hoteldto.ListAdditionalRoomsResponse, error)
	ListAllBedTypes(ctx context.Context, req *hoteldto.ListAllBedTypesRequest) (*hoteldto.ListAllBedTypesResponse, error)
	DetailHotel(ctx context.Context, hotelID uint) (*hoteldto.DetailHotelResponse, error)
	DetailHotelForAgent(ctx context.Context, hotelID uint) (*hoteldto.DetailHotelForAgentResponse, error)
	RemoveHotel(ctx context.Context, hotelID uint) error
	RemoveRoomType(ctx context.Context, roomTypeID uint) error
	ListRoomAvailable(ctx context.Context, req *hoteldto.ListRoomAvailableRequest) (*hoteldto.ListRoomAvailableResponse, error)
	UpdateRoomAvailable(ctx context.Context, req *hoteldto.UpdateRoomAvailableRequest) error
	ListProvinces(ctx context.Context, request *hoteldto.ListProvincesRequest) (*hoteldto.ListProvincesResponse, error)
	AddRoomType(ctx context.Context, hotelID uint, req *hoteldto.AddRoomTypeRequest) error
	UpdateStatus(ctx context.Context, req *hoteldto.UpdateStatusRequest) error
	ListStatusHotel(ctx context.Context) (*hoteldto.ListStatusHotelResponse, error)
	UpdateHotel(ctx context.Context, req *hoteldto.UpdateHotelRequest) error
	UpdateRoomType(ctx context.Context, req *hoteldto.UpdateRoomTypeRequest) error
}

type HotelRepository interface {
	GetHotels(ctx context.Context, filter filter.HotelFilter) ([]entity.Hotel, int64, error)
	GetHotelsForAgent(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.CustomHotel, int64, error)
	GetFilterDistricts(ctx context.Context, filter filter.HotelFilterForAgent) ([]string, error)
	GetFilterPricing(ctx context.Context, filter filter.HotelFilterForAgent) (*entity.FilterRangePrice, error)
	GetFilterRatings(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterRatingHotel, error)
	GetFilterBedTypes(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterBedTypeHotel, error)
	GetFilterTotalBedrooms(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterTotalBedroom, error)
	GetRoomTypeByHotelID(ctx context.Context, hotelID uint) ([]entity.RoomType, error)
	GetBedTypeByRoomTypeID(ctx context.Context, roomTypeID uint) ([]entity.BedType, error)
	CreateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error)
	AttachPhotosHotel(ctx context.Context, hotelID uint, photoURLs []string) error
	AttachFacilities(ctx context.Context, hotelID uint, facilityNames []string) error
	AttachNearbyPlaces(ctx context.Context, hotelID uint, np []hoteldto.NearbyPlace) error
	CreateRoomType(ctx context.Context, roomType *entity.RoomType) (*entity.RoomType, error)
	AttachPhotosRoomType(ctx context.Context, roomTypeID uint, photoURLs []string) error
	AttachRoomAdditions(ctx context.Context, roomTypeID uint, additionals []entity.CustomRoomAdditional) error
	AttachBedTypesToRoomType(ctx context.Context, roomTypeID uint, bedTypeNames []string) error
	CreateRoomPrice(ctx context.Context, roomTypeID uint, dto *entity.CustomBreakfast, isBreakfast bool) error
	GetFacilities(ctx context.Context, filter *filter.DefaultFilter) ([]string, int64, error)
	GetAdditionalRooms(ctx context.Context, filter *filter.DefaultFilter) ([]string, int64, error)
	GetBedTypes(ctx context.Context, filter *filter.DefaultFilter) ([]string, int64, error)
	GetHotelByID(ctx context.Context, hotelID uint, scope string) (*entity.Hotel, error)
	DeleteRoomType(ctx context.Context, roomTypeID uint) error
	DeleteHotel(ctx context.Context, hotelID uint) error
	GetRoomUnavailableByRoomTypeIDs(ctx context.Context, roomTypeIDs []uint, month time.Time) ([]entity.RoomUnavailable, error)
	DeleteRoomUnavailable(ctx context.Context, roomTypeID uint, month time.Time) error
	InsertRoomUnavailable(ctx context.Context, roomTypeID uint, unavailableDates []time.Time) error
	GetProvinces(ctx context.Context, filter *filter.DefaultFilter) ([]string, int64, error)
	GetRoomPriceByID(ctx context.Context, id uint) (*entity.RoomPrice, error)
	GetRoomTypeAdditionalsByIDs(ctx context.Context, ids []uint) ([]entity.RoomTypeAdditional, error)
	UpdateStatus(ctx context.Context, hotelID uint, statusID uint) error
	UpdateHotel(ctx context.Context, hotel *entity.Hotel) error
	GetRoomTypeByID(ctx context.Context, roomTypeID uint) (*entity.RoomType, error)
	UpdateRoomType(ctx context.Context, roomType *entity.RoomType) error
}

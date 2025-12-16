package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
)

type BookingUsecase interface {
	AddToCart(ctx context.Context, req *bookingdto.AddToCartRequest) error
	ListCart(ctx context.Context) (*bookingdto.ListCartResponse, error)
	RemoveFromCart(ctx context.Context, bookingDetailID uint) error
	CheckOutCart(ctx context.Context) (*bookingdto.CheckOutCartResponse, error)
	ListBookingHistory(ctx context.Context, req *bookingdto.ListBookingHistoryRequest) (*bookingdto.ListBookingHistoryResponse, error)
	ListBookings(ctx context.Context, req *bookingdto.ListBookingsRequest) (*bookingdto.ListBookingsResponse, error)
	ListBookingLog(ctx context.Context, req *bookingdto.ListBookingLogRequest) (*bookingdto.ListBookingLogResponse, error)
	UploadReceipt(ctx context.Context, req *bookingdto.UploadReceiptRequest) error
	UpdateStatusBooking(ctx context.Context, req *bookingdto.UpdateStatusRequest, scope string) error
	ListStatusBooking(ctx context.Context) (*bookingdto.ListStatusBookingResponse, error)
	ListStatusPayment(ctx context.Context) (*bookingdto.ListStatusPaymentResponse, error)
	ListBookingIDs(ctx context.Context, req *bookingdto.ListBookingIDsRequest) (*bookingdto.ListBookingIDsResponse, error)
	ListSubBookingIDs(ctx context.Context, req *bookingdto.ListSubBookingIDsRequest) (*bookingdto.ListSubBookingIDsResponse, error)
	AddGuestsToCart(ctx context.Context, req *bookingdto.AddGuestsToCartRequest) error
	RemoveGuestsFromCart(ctx context.Context, req *bookingdto.RemoveGuestsFromCartRequest) error
	AddGuestToSubCart(ctx context.Context, req *bookingdto.AddGuestToSubCartRequest) error
	CancelBooking(ctx context.Context, req *bookingdto.CancelBookingRequest) error
}

type BookingRepository interface {
	GetOrCreateCartID(ctx context.Context, agentID uint) (uint, error)
	CreateBookingDetail(ctx context.Context, detail *entity.BookingDetail) ([]uint, error)
	CreateBookingDetailAdditional(ctx context.Context, add *entity.BookingDetailAdditional) error
	GetCartBooking(ctx context.Context, agentID uint) (*entity.Booking, error)
	DeleteCartBooking(ctx context.Context, agentID uint, bookingDetailID uint) error
	UpdateBookingGuests(ctx context.Context, bookingID uint, guests []string) error
	UpdateBookingDetailGuest(ctx context.Context, detailID uint, guest string) error
	UpdateBookingStatus(ctx context.Context, bookingID uint, statusID uint) error
	GetBookings(ctx context.Context, filter *filter.BookingFilter) ([]entity.Booking, int64, error)
	UpdateBookingReceipt(ctx context.Context, bookingDetailID []uint, receiptURL string) error
	GetBookingByID(ctx context.Context, bookingID uint) (*entity.Booking, error)
	UpdateBookingDetailStatusBooking(ctx context.Context, bookingDetailID []uint, statusID uint) ([]entity.BookingDetail, []string, error)
	UpdateBookingDetailStatusPayment(ctx context.Context, bookingDetailID []uint, statusID uint) error
	GetBookingByCode(ctx context.Context, code string) (*entity.Booking, error)
	GetSubBookingByCode(ctx context.Context, code string) (*entity.BookingDetail, error)
	GetBookingIDs(ctx context.Context, agentID uint, filter *filter.DefaultFilter) ([]string, int64, error)
	GetSubBookingIDs(ctx context.Context, agentID uint, bookingCode string) ([]string, error)
	AddGuestsToCart(ctx context.Context, agentID uint, bookingID uint, guests []bookingdto.GuestInfo) error
	RemoveGuestsFromCart(ctx context.Context, agentID uint, bookingID uint, guests []bookingdto.GuestInfo) error
	AddGuestToSubCart(ctx context.Context, agentID uint, bookingDetailID uint, guest string) error
	CancelBooking(ctx context.Context, agentID uint, subBookingID string) (*entity.BookingDetail, error)
	CreateInvoice(ctx context.Context, invoices []entity.Invoice) error
	GenerateCode(ctx context.Context, keyRedis string, prefixCode string) (string, error)
	GetBookingDetailIDsByBookingCode(ctx context.Context, bookingCode string) ([]uint, error)
	GetIDBySubBookingID(ctx context.Context, subBookingID string) (uint, error)
	GetListBookingLog(ctx context.Context, filter *filter.BookingFilter) ([]entity.BookingDetail, int64, error)
	UpdateDetailBookingDetail(ctx context.Context, bookingDetailID uint, room *entity.DetailRoom, promo *entity.DetailPromo, price float64, additionals []entity.BookingDetailAdditional) error
	GetBookingGuests(ctx context.Context, bookingID uint) ([]model.BookingGuest, error)
}

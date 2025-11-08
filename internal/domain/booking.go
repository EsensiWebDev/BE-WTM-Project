package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
)

type BookingUsecase interface {
	AddToCart(ctx context.Context, req *bookingdto.AddToCartRequest) error
	ListCart(ctx context.Context) (*bookingdto.ListCartResponse, error)
	RemoveFromCart(ctx context.Context, bookingDetailID uint) error
	CheckOutCart(ctx context.Context, req *bookingdto.CheckOutCartRequest) error
	ListBookingHistory(ctx context.Context, req *bookingdto.ListBookingHistoryRequest) (*bookingdto.ListBookingHistoryResponse, error)
	ListBookings(ctx context.Context, req *bookingdto.ListBookingsRequest) (*bookingdto.ListBookingsResponse, error)
	ListBookingLog(ctx context.Context, req *bookingdto.ListBookingLogRequest) (*bookingdto.ListBookingLogResponse, error)
	UploadReceipt(ctx context.Context, req *bookingdto.UploadReceiptRequest) error
	UpdateStatusBooking(ctx context.Context, req *bookingdto.UpdateStatusBookingRequest) error
	ListStatusBooking(ctx context.Context) (*bookingdto.ListStatusBookingResponse, error)
	ListStatusPayment(ctx context.Context) (*bookingdto.ListStatusPaymentResponse, error)
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
	UpdateBookingDetailStatus(ctx context.Context, bookingDetailID []uint, statusID uint) ([]entity.BookingDetail, error)
	GetBookingByCode(ctx context.Context, code string) (*entity.Booking, error)
	GetSubBookingByCode(ctx context.Context, code string) (*entity.BookingDetail, error)
}

package booking_usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/pkg/constant"
)

type BookingUsecase struct {
	bookingRepo domain.BookingRepository
	hotelRepo   domain.HotelRepository
	promoRepo   domain.PromoRepository
	middleware  domain.Middleware
	dbTrx       domain.DatabaseTransaction
	fileStorage domain.StorageClient
	config      *config.Config
	emailRepo   domain.EmailRepository
	emailSender domain.EmailSender
	userRepo    domain.UserRepository
}

func NewBookingUsecase(bookingRepo domain.BookingRepository, hotelRepo domain.HotelRepository, promoRepo domain.PromoRepository, middleware domain.Middleware, dbTrx domain.DatabaseTransaction, fileStorage domain.StorageClient, config *config.Config, emailRepo domain.EmailRepository, emailSender domain.EmailSender, userRepo domain.UserRepository) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo: bookingRepo,
		hotelRepo:   hotelRepo,
		promoRepo:   promoRepo,
		middleware:  middleware,
		dbTrx:       dbTrx,
		fileStorage: fileStorage,
		config:      config,
		emailRepo:   emailRepo,
		emailSender: emailSender,
		userRepo:    userRepo,
	}
}

func (bu *BookingUsecase) uploadFile(ctx context.Context, file *multipart.FileHeader, prefix string, bookingID uint) (string, error) {
	if file == nil {
		return "", nil
	}
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("failed to close file:", err.Error())
		}
	}(f)

	bucketName := fmt.Sprintf("%s-%s", constant.ConstBooking, constant.ConstPrivate)
	filename := fmt.Sprintf("%s_%d_%d%s", prefix, bookingID, time.Now().Unix(), filepath.Ext(file.Filename))
	return bu.fileStorage.UploadFile(ctx, f, file, bucketName, filename)
}

func (bu *BookingUsecase) summaryStatus(statuses []string, types string) string {
	if len(statuses) == 0 {
		return "No Status"
	}

	var priority []string
	switch types {
	case constant.ConstBooking:
		priority = []string{
			constant.StatusBookingRejected,
			constant.StatusBookingWaitingApproval,
			constant.StatusBookingConfirmed,
			constant.StatusBookingCanceled,
		}
	case constant.ConstPayment:
		priority = []string{
			constant.StatusPaymentUnpaid,
			constant.StatusPaymentPaid,
		}
	}

	// cek sesuai urutan prioritas
	for _, p := range priority {
		count := 0
		for _, s := range statuses {
			if s == p {
				count++
			}
		}
		if count > 0 {
			if count == len(statuses) {
				return p
			}
			return fmt.Sprintf("%d of %d %s", count, len(statuses), p)
		}
	}

	return "Unknown Status"

}

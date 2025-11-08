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
}

func NewBookingUsecase(bookingRepo domain.BookingRepository, hotelRepo domain.HotelRepository, promoRepo domain.PromoRepository, middleware domain.Middleware, dbTrx domain.DatabaseTransaction, fileStorage domain.StorageClient, config *config.Config, emailRepo domain.EmailRepository, emailSender domain.EmailSender) *BookingUsecase {
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
	}
}

func (bu *BookingUsecase) uploadFile(ctx context.Context, file *multipart.FileHeader, prefix string, bookingDetailID uint) (string, error) {
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
	filename := fmt.Sprintf("%s_%d_%d%s", prefix, bookingDetailID, time.Now().Unix(), filepath.Ext(file.Filename))
	return bu.fileStorage.UploadFile(ctx, f, file, bucketName, filename)
}

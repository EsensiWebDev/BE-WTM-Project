package email_usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"wtm-backend/internal/domain"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

type EmailUsecase struct {
	emailRepo   domain.EmailRepository
	emailSender domain.EmailSender
	bookingRepo domain.BookingRepository
	fileStorage domain.StorageClient
}

func NewEmailUsecase(emailRepo domain.EmailRepository, emailSender domain.EmailSender, bookingRepo domain.BookingRepository) *EmailUsecase {
	return &EmailUsecase{
		emailRepo:   emailRepo,
		emailSender: emailSender,
		bookingRepo: bookingRepo,
	}
}

func (eu *EmailUsecase) uploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
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
			logger.Error(ctx, "failed to close file", err.Error())
		}
	}(f)

	bucketName := fmt.Sprintf("%s-%s", constant.ConstEmail, constant.ConstPublic)
	filename := fmt.Sprintf("%s-%d%s", constant.ConstSignature, time.Now().Unix(), filepath.Ext(file.Filename))
	return eu.fileStorage.UploadFile(ctx, f, file, bucketName, filename)
}

package hotel_usecase

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"mime/multipart"
	"path"
	"path/filepath"
	"sync"
	"time"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

type HotelUsecase struct {
	hotelRepo     domain.HotelRepository
	fileStorage   domain.StorageClient
	dbTransaction domain.DatabaseTransaction
	config        *config.Config
}

func NewHotelUsecase(hotelRepo domain.HotelRepository, fileStorage domain.StorageClient, dbTrx domain.DatabaseTransaction, config *config.Config) *HotelUsecase {
	return &HotelUsecase{
		hotelRepo:     hotelRepo,
		fileStorage:   fileStorage,
		dbTransaction: dbTrx,
		config:        config,
	}
}

func (hu *HotelUsecase) uploadMultiple(
	ctx context.Context,
	files []*multipart.FileHeader,
	typeAccess string,
	prefixParts ...string,
) ([]string, error) {
	var (
		urls []string
		mu   sync.Mutex
	)

	bucketName := fmt.Sprintf("%s-%s", constant.ConstHotel, typeAccess)
	prefix := path.Join(prefixParts...)

	g, ctx := errgroup.WithContext(ctx)

	for i, fh := range files {
		i, fh := i, fh // avoid closure capture bug

		g.Go(func() error {
			file, err := fh.Open()
			if err != nil {
				logger.Error(ctx, "failed to open file", err.Error())
				return fmt.Errorf("cannot open file %s: %s", fh.Filename, err.Error())
			}
			defer func(file multipart.File) {
				err := file.Close()
				if err != nil {
					logger.Error(ctx, "failed to close file", err.Error())
				}
			}(file)

			filename := fmt.Sprintf("%s_%d_%d%s",
				prefix, time.Now().UnixNano(), i, filepath.Ext(fh.Filename))

			url, err := hu.fileStorage.UploadFile(ctx, file, fh, bucketName, filename)
			if err != nil {
				logger.Error(ctx, "upload error", err.Error())
				return fmt.Errorf("upload failed for %s: %s", fh.Filename, err.Error())
			}

			mu.Lock()
			urls = append(urls, url)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logger.Error(ctx, "Error uploading files", err.Error())
		return nil, err
	}
	return urls, nil
}

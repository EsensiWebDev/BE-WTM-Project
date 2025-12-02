package banner_usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) UpsertBanner(ctx context.Context, req *bannerdto.UpsertBannerRequest, reqId *bannerdto.DetailBannerRequest) error {

	var err error
	banner := &entity.Banner{
		Title:       req.Title,
		Description: req.Description,
	}

	if reqId == nil {

		banner, err = bu.bannerRepo.CreateBanner(ctx, banner)
		if err != nil {
			logger.Error(ctx, "Error creating banner", err.Error())
			return err
		}

	} else {

		banner, err = bu.bannerRepo.GetBannerByExternalID(ctx, reqId.BannerID)
		if err != nil {
			logger.Error(ctx, "Error getting banner by Id", err.Error())
			return err
		}

		if banner == nil {
			logger.Warn(ctx, "Banner not found", "Id", reqId.BannerID)
			return fmt.Errorf("banner with Id %s not found", reqId.BannerID)
		}

		if banner.Title == req.Title && banner.Description == req.Description && req.Image == nil {
			logger.Warn(ctx, "Banner not updated", "Id", reqId.BannerID)
			return errors.New("no changes detected on banner")
		}

		if banner.Title != req.Title {
			banner.Title = req.Title
		}

		if banner.Description != req.Description {
			banner.Description = req.Description
		}

	}

	if req.Image != nil && req.Image.Size > 0 {
		ImageUrl, err := bu.uploadFile(ctx, banner.ID, req.Image)
		if err != nil {
			logger.Error(ctx, "Error uploading banner image", err.Error())
			return err
		}

		banner.ImageURL = ImageUrl
	}

	if err = bu.bannerRepo.UpdateBanner(ctx, banner); err != nil {
		logger.Error(ctx, "Error updating banner", err.Error())
		return err
	}

	return nil

}

func (bu *BannerUsecase) uploadFile(ctx context.Context, bannerID uint, file *multipart.FileHeader) (string, error) {
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

	bucketName := fmt.Sprintf("%s-%s", constant.ConstBanner, constant.ConstPublic)
	filename := fmt.Sprintf("%d_%d%s", bannerID, time.Now().Unix(), filepath.Ext(file.Filename))
	return bu.fileStorage.UploadFile(ctx, f, file, bucketName, filename)
}

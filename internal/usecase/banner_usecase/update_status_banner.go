package banner_usecase

import (
	"context"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) UpdateStatusBanner(ctx context.Context, req *bannerdto.UpdateStatusBannerRequest) error {
	if err := bu.bannerRepo.UpdateStatusBanner(ctx, req.ID, req.Status); err != nil {
		logger.Error(ctx, "Error updating banner status:", err.Error())
		return err
	}

	return nil
}

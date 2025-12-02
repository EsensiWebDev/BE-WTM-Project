package banner_usecase

import (
	"context"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) UpdateOrderBanner(ctx context.Context, req *bannerdto.UpdateOrderBannerRequest) error {
	if err := bu.bannerRepo.UpdateOrderBanner(ctx, req.ID, req.Order); err != nil {
		logger.Error(ctx, "update order banner failed", err.Error())
		return err
	}
	return nil
}

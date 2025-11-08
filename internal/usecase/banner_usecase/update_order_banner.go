package banner_usecase

import (
	"context"
	"errors"
	"wtm-backend/internal/dto/bannerdto"
)

func (bu *BannerUsecase) UpdateOrderBanner(ctx context.Context, req *bannerdto.UpdateOrderBannerRequest) error {
	return bu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {
		for _, banner := range req.Data {
			if banner.ID > 0 && banner.Order > 0 {
				if err := bu.bannerRepo.UpdateOrderBanner(txCtx, banner.ID, banner.Order); err != nil {
					return err
				}
			} else {
				return errors.New("invalid order banner")
			}
		}
		return nil
	})
}

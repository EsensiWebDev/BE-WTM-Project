package promo_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/pkg/logger"
)

func (pu *PromoUsecase) SetStatusPromo(ctx context.Context, req *promodto.SetStatusPromoRequest) error {

	selectedFields := []string{"id", "is_active"}

	promo, err := pu.promoRepo.GetPromoByID(ctx, req.PromoID, selectedFields)
	if err != nil {
		logger.Error(ctx, "Error getting promo by Id", "error", err, "promoID", req.PromoID)
		return err
	}

	if promo == nil {
		logger.Error(ctx, "Promo not found", "promoID", req.PromoID)
		return fmt.Errorf("promo not found")
	}

	if promo.IsActive == req.IsActive {
		logger.Info(ctx, "Promo status is already set", "promoID", req.PromoID, "isActive", req.IsActive)
		return nil
	}

	if err := pu.promoRepo.UpdatePromoStatus(ctx, req.PromoID, req.IsActive); err != nil {
		logger.Error(ctx, "Error updating promo status", "error", err, "promoID", req.PromoID, "isActive", req.IsActive)
		return err
	}

	return nil

}

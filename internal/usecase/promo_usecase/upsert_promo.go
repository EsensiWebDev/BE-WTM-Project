package promo_usecase

import (
	"context"
	"strconv"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pu *PromoUsecase) UpsertPromo(ctx context.Context, req *promodto.UpsertPromoRequest, promoID *uint) error {
	return pu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {
		startDate, err := utils.ParseRFC3339ToUTC(req.StartDate)
		if err != nil {
			logger.Error(ctx, "Error parsing start date", err.Error())
			return err
		}

		endDate, err := utils.ParseRFC3339ToUTC(req.EndDate)
		if err != nil {
			logger.Error(ctx, "Error parsing end date", err.Error())
			return err
		}

		var promoRoomTypes []entity.PromoRoomType
		if len(req.RoomTypes) > 0 {
			for _, roomType := range req.RoomTypes {
				promoRoomTypes = append(promoRoomTypes, entity.PromoRoomType{
					RoomTypeID:  roomType.RoomTypeID,
					TotalNights: roomType.TotalNight,
				})
			}
		}

		var detail entity.PromoDetail

		switch req.PromoTypeID {
		case constant.PromoTypeDiscountID:
			discount, err := strconv.ParseFloat(req.Detail, 64)
			if err != nil {
				logger.Error(ctx, "Error parsing discount id", err.Error())
			}
			if discount > 0 {
				detail = entity.PromoDetail{
					DiscountPercentage: discount,
				}
			}
		case constant.PromoTypeFixedPriceID:
			fixedPrice, err := strconv.ParseFloat(req.Detail, 64)
			if err != nil {
				logger.Error(ctx, "Error parsing fixed price", err.Error())
			}
			if fixedPrice > 0 {
				detail = entity.PromoDetail{
					FixedPrice: fixedPrice,
				}
			}
		case constant.PromoTypeRoomUpgradeID:
			roomUpgradeID, err := utils.StringToUint(req.Detail)
			if err != nil {
				logger.Error(ctx, "Error parsing room upgrade Id", err.Error())
			}

			if roomUpgradeID > 0 {
				detail = entity.PromoDetail{
					UpgradedToID: roomUpgradeID,
				}
			}
		case constant.PromoTypeBenefitID:
			if req.Detail != "" {
				detail = entity.PromoDetail{
					BenefitNote: req.Detail,
				}
			}
		}

		promo := &entity.Promo{
			Name:           req.PromoName,
			Description:    req.Description,
			Code:           req.PromoCode,
			PromoTypeID:    req.PromoTypeID,
			Detail:         detail,
			IsActive:       false,
			StartDate:      &startDate,
			EndDate:        &endDate,
			PromoRoomTypes: promoRoomTypes,
		}

		if promoID != nil && *promoID > 0 {
			promo.ID = *promoID
			err = pu.promoRepo.UpdatePromo(txCtx, promo)
			if err != nil {
				logger.Error(ctx, "Error updating promo", err.Error())
				return err
			}
		} else {
			err = pu.promoRepo.CreatePromo(ctx, promo)
			if err != nil {
				logger.Error(ctx, "Error creating promo", err.Error())
				return err
			}
		}

		return nil
	})
}

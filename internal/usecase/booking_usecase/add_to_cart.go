package booking_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) AddToCart(ctx context.Context, req *bookingdto.AddToCartRequest) error {
	return bu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {

		// Get agent Id from context
		userCtx, err := bu.middleware.GenerateUserFromContext(txCtx)
		if err != nil {
			logger.Error(ctx, "failed to get user from context", err.Error())
			return fmt.Errorf("failed to get user from context: %s", err.Error())
		}

		if userCtx == nil {
			logger.Error(ctx, "user context is nil")
			return fmt.Errorf("user context is nil")
		}

		agentID := userCtx.ID

		//Get RoomPrice
		roomPrice, err := bu.hotelRepo.GetRoomPriceByID(txCtx, req.RoomPriceID)
		if err != nil {
			logger.Error(ctx, "failed to get room price by id", err.Error())
			return fmt.Errorf("room price not found: %s", err.Error())
		}

		//Get Promo (optional)
		var promo *entity.Promo
		if req.PromoID > 0 {
			promo, err = bu.promoRepo.GetPromoByID(txCtx, req.PromoID, nil)
			if err != nil {
				logger.Error(ctx, "failed to get promo by id", err.Error())
				return fmt.Errorf("promo not found: %s", err.Error())
			}
		}

		//Get Additionals
		additionals, err := bu.hotelRepo.GetRoomTypeAdditionalsByIDs(txCtx, req.RoomTypeAdditionalIDs)
		if err != nil {
			logger.Error(ctx, "failed to get room type additional", err.Error())
			return fmt.Errorf("additionals not found: %s", err.Error())
		}

		// Create BookingDetail

		bookingID, err := bu.bookingRepo.GetOrCreateCartID(txCtx, agentID)
		if err != nil {
			logger.Error(ctx, "failed to get or create cart Id", err.Error())
			return fmt.Errorf("failed to get or create cart Id: %s", err.Error())
		}

		if bookingID == 0 {
			logger.Error(ctx, "booking Id is zero, cart creation failed")
			return fmt.Errorf("failed to create booking cart")
		}

		checkInDate, err := time.Parse(time.DateOnly, req.CheckInDate)
		if err != nil {
			logger.Error(ctx, "failed to parse check-in date", err.Error())
			return fmt.Errorf("invalid RFC3339 date: %s", err.Error())
		}

		checkOutDate, err := time.Parse(time.DateOnly, req.CheckOutDate)
		if err != nil {
			logger.Error(ctx, "failed to parse check-out date", err.Error())
			return fmt.Errorf("invalid check-out date: %s", err.Error())
		}

		//var photoUrl string
		//for _, photo := range roomPrice.RoomType.Hotel.Photos {
		//	if photo != "" {
		//		bucketName := fmt.Sprintf("%s-%s", constant.ConstHotel, constant.ConstPublic)
		//		photoUrl, err = bu.fileStorage.GetFile(ctx, bucketName, photo)
		//		if err != nil {
		//			logger.Error(ctx, "Error getting banner image", err.Error())
		//			continue
		//		}
		//		break
		//	}
		//}

		//var cancelledDate time.Time
		//cancellationPeriod := roomPrice.RoomType.Hotel.CancellationPeriod
		//if cancellationPeriod > 0 {
		//	cancelledDate = checkInDate.AddDate(0, 0, -cancellationPeriod)
		//}
		//detailRoom := entity.DetailRoom{
		//	OriPrice:      roomPrice.Price,
		//	HotelName:     roomPrice.RoomType.Hotel.Name,
		//	HotelRating:   roomPrice.RoomType.Hotel.Rating,
		//	HotelPhoto:    photoUrl,
		//	RoomTypeName:  roomPrice.RoomType.Name,
		//	IsBreakfast:   roomPrice.IsBreakfast,
		//	CancelledDate: cancelledDate.Format(time.DateOnly),
		//	IsAPI:         roomPrice.RoomType.Hotel.IsAPI,
		//	Capacity:      roomPrice.RoomType.MaxOccupancy,
		//}

		// Hitung jumlah malam
		nights := int(checkOutDate.Sub(checkInDate).Hours() / 24)
		if nights <= 0 {
			logger.Error(ctx, "invalid stay duration")
			return fmt.Errorf("check-out date must be after check-in date")
		}

		if promo != nil {
			if promo.Duration > nights {
				logger.Error(ctx, "This promo is not valid for the selected stay duration")
				return fmt.Errorf("this promo is not valid for the selected stay duration")
			}
		}

		detailBooking := &entity.BookingDetail{
			BookingID:   bookingID,
			RoomPriceID: roomPrice.ID,
			//RoomTypeID:      roomPrice.RoomType.ID,
			CheckInDate:  checkInDate,
			CheckOutDate: checkOutDate,
			Quantity:     req.Quantity,
			//DetailRooms:     detailRoom,
			StatusBookingID: constant.StatusBookingWaitingApprovalID,
			StatusPaymentID: constant.StatusPaymentUnpaidID,
		}

		if promo != nil {
			detailBooking.PromoID = &promo.ID
		}

		//basePrice := roomPrice.Price
		//var finalRoomPrice float64
		//if promo != nil {
		//	switch promo.PromoTypeName {
		//	case constant.PromoTypeFixedPrice:
		//		finalRoomPrice = promo.Detail.FixedPrice
		//		if promo.Duration > nights {
		//			finalRoomPrice += float64(nights-promo.Duration) * basePrice
		//		}
		//	case constant.PromoTypeDiscount:
		//		finalRoomPrice = (100 - promo.Detail.DiscountPercentage) / 100 * basePrice * float64(nights)
		//		if promo.Duration > nights {
		//			finalRoomPrice += float64(nights-promo.Duration) * basePrice
		//		}
		//	default:
		//		finalRoomPrice = basePrice * float64(nights)
		//	}
		//	detailPromo, err := bu.generateDetailPromo(promo)
		//	if err != nil {
		//		logger.Error(ctx, "failed to generate detail promo", err.Error())
		//	}
		//	detailBooking.PromoID = &promo.ID
		//	detailBooking.DetailPromos = detailPromo
		//} else {
		//	finalRoomPrice = basePrice * float64(nights)
		//}
		//
		//detailBooking.Price = finalRoomPrice

		bookingDetailIds, err := bu.bookingRepo.CreateBookingDetail(txCtx, detailBooking)
		if err != nil {
			logger.Error(ctx, "failed to create booking detail", err.Error())
			return fmt.Errorf("failed to create booking detail: %s", err.Error())
		}

		// 6. Create BookingDetailAdditionals
		for _, add := range additionals {
			additional := &entity.BookingDetailAdditional{
				BookingDetailIDs:     bookingDetailIds,
				RoomTypeAdditionalID: add.ID,
				Category:             add.Category,
				Price:                add.Price,
				Pax:                  add.Pax,
				IsRequired:           add.IsRequired,
				NameAdditional:       add.RoomAdditional.Name,
			}
			if err := bu.bookingRepo.CreateBookingDetailAdditional(txCtx, additional); err != nil {
				return fmt.Errorf("failed to create additional: %s", err.Error())
			}
		}

		return nil
	})
}

func (bu *BookingUsecase) generateDetailPromo(promo *entity.Promo) (entity.DetailPromo, error) {

	detailPromo := entity.DetailPromo{
		Name:            promo.Name,
		PromoCode:       promo.Code,
		Type:            promo.PromoTypeName,
		DiscountPercent: promo.Detail.DiscountPercentage,
		FixedPrice:      promo.Detail.FixedPrice,
		UpgradedToID:    promo.Detail.UpgradedToID,
		BenefitNote:     promo.Detail.BenefitNote,
	}

	return detailPromo, nil
}

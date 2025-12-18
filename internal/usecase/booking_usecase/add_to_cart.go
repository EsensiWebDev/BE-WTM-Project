package booking_usecase

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/currency"
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

		checkInDate, err := time.Parse(time.DateOnly, req.CheckInDate)
		if err != nil {
			logger.Error(ctx, "failed to parse check-in date", err.Error())
			return fmt.Errorf("invalid check-in date: %s", err.Error())
		}

		checkOutDate, err := time.Parse(time.DateOnly, req.CheckOutDate)
		if err != nil {
			logger.Error(ctx, "failed to parse check-out date", err.Error())
			return fmt.Errorf("invalid check-out date: %s", err.Error())
		}

		// Validate booking limit per booking if set
		if roomPrice.RoomType.BookingLimitPerBooking != nil {
			bookingLimit := *roomPrice.RoomType.BookingLimitPerBooking
			if bookingLimit > 0 {
				// Get existing cart to check current quantities for this room type
				existingCart, err := bu.bookingRepo.GetCartBooking(txCtx, agentID)
				if err != nil {
					logger.Error(ctx, "failed to get cart booking for limit validation", err.Error())
					return fmt.Errorf("failed to validate booking limit: %s", err.Error())
				}

				// Sum up quantities for the same room type and overlapping date ranges
				totalQuantity := req.Quantity
				if existingCart != nil {
					for _, detail := range existingCart.BookingDetails {
						// Check if it's the same room type (using RoomType.ID since RoomType is preloaded)
						if detail.RoomPrice.RoomType.ID == roomPrice.RoomType.ID {
							// Check if date ranges overlap
							// Two date ranges overlap if: checkInDate < detail.CheckOutDate && checkOutDate > detail.CheckInDate
							if checkInDate.Before(detail.CheckOutDate) && checkOutDate.After(detail.CheckInDate) {
								totalQuantity += detail.Quantity
							}
						}
					}
				}

				if totalQuantity > bookingLimit {
					logger.Error(ctx, fmt.Sprintf("Booking limit exceeded for room type %s. Limit: %d, Requested total: %d", roomPrice.RoomType.Name, bookingLimit, totalQuantity))
					return fmt.Errorf("booking limit exceeded: maximum %d rooms allowed per booking for %s, but %d rooms requested", bookingLimit, roomPrice.RoomType.Name, totalQuantity)
				}
			}
		}

		// Validate bed type if provided
		if req.BedType != "" {
			bedTypeValid := false
			for _, bedTypeName := range roomPrice.RoomType.BedTypeNames {
				if bedTypeName == req.BedType {
					bedTypeValid = true
					break
				}
			}
			if !bedTypeValid {
				logger.Error(ctx, "Invalid bed type selected", req.BedType)
				return fmt.Errorf("invalid bed type selected: %s", req.BedType)
			}
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

		//Get Other Preferences
		var preferences []entity.RoomTypePreference
		if len(req.OtherPreferenceIDs) > 0 {
			preferences, err = bu.hotelRepo.GetRoomTypePreferencesByIDs(txCtx, req.OtherPreferenceIDs)
			if err != nil {
				logger.Error(ctx, "failed to get room type preferences", err.Error())
				return fmt.Errorf("preferences not found: %s", err.Error())
			}
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

		// Join selected "Other Preferences" into a comma-separated string snapshot
		var otherPrefsJoined string
		if len(preferences) > 0 {
			var preferenceNames []string
			for _, pref := range preferences {
				preferenceNames = append(preferenceNames, pref.OtherPreference.Name)
			}
			otherPrefsJoined = strings.Join(preferenceNames, ", ")
		}

		// Trim additional notes (admin-only field)
		additionalNotes := strings.TrimSpace(req.AdditionalNotes)

		// Get agent's currency preference
		user, err := bu.userRepo.GetUserByID(txCtx, agentID)
		if err != nil {
			logger.Error(ctx, "failed to get user for currency", err.Error())
			return fmt.Errorf("failed to get user: %s", err.Error())
		}
		agentCurrency := "IDR" // Default fallback
		if user != nil && user.Currency != "" {
			agentCurrency = user.Currency
		}

		detailBooking := &entity.BookingDetail{
			BookingID:        bookingID,
			RoomPriceID:      roomPrice.ID,
			CheckInDate:      checkInDate,
			CheckOutDate:     checkOutDate,
			Quantity:         req.Quantity,
			StatusBookingID:  constant.StatusBookingWaitingApprovalID,
			StatusPaymentID:  constant.StatusPaymentUnpaidID,
			BedType:          req.BedType,
			OtherPreferences: otherPrefsJoined,
			AdditionalNotes:  additionalNotes,
			Currency:         agentCurrency, // Store agent's currency at booking time
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
			var price *float64

			// Convert price to agent's currency if it's a price-based additional
			if add.Category == constant.AdditionalServiceCategoryPrice {
				if len(add.Prices) > 0 {
					// Use Prices map for multi-currency support
					if convertedPrice, _, _ := currency.GetPriceForCurrency(add.Prices, agentCurrency); convertedPrice > 0 {
						price = &convertedPrice
					} else if add.Price != nil && *add.Price > 0 {
						// Fallback to Price field if Prices not available (backward compatibility)
						price = add.Price
					}
				} else if add.Price != nil && *add.Price > 0 {
					// Backward compatibility: use Price field if Prices not set
					price = add.Price
				}
			}

			additional := &entity.BookingDetailAdditional{
				BookingDetailIDs:     bookingDetailIds,
				RoomTypeAdditionalID: add.ID,
				Category:             add.Category,
				Price:                price,
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
		Description:     promo.Description,
		PromoTypeID:     promo.PromoTypeID,
		DiscountPercent: promo.Detail.DiscountPercentage,
		FixedPrice:      promo.Detail.FixedPrice,
		Prices:          promo.Detail.Prices,
		UpgradedToID:    promo.Detail.UpgradedToID,
		BenefitNote:     promo.Detail.BenefitNote,
	}

	return detailPromo, nil
}

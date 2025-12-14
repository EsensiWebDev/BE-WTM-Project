package booking_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListCart(ctx context.Context) (*bookingdto.ListCartResponse, error) {
	// Get agent Id from context
	userCtx, err := bu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return nil, fmt.Errorf("failed to get user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return nil, fmt.Errorf("user context is nil")
	}

	agentID := userCtx.ID

	cart, err := bu.bookingRepo.GetCartBooking(ctx, agentID)
	if err != nil {
		logger.Error(ctx, "failed to get cart booking", err.Error())
		return nil, fmt.Errorf("failed to get cart booking: %s", err.Error())
	}

	result := &bookingdto.ListCartResponse{}
	if cart != nil {
		result.ID = cart.ID
		var details []bookingdto.CartDetail
		var grandTotal float64

		for _, detail := range cart.BookingDetails {
			var additionals []bookingdto.CartDetailAdditional
			var totalAdditional float64

			// Map additional services with new structure (Category, Price/Pax, IsRequired)
			// These fields are now stored directly in BookingDetailAdditional
			for _, detailAdditional := range detail.BookingDetailsAdditional {
				cartAdditional := bookingdto.CartDetailAdditional{
					Name:       detailAdditional.NameAdditional,
					Category:   detailAdditional.Category,
					Price:      detailAdditional.Price,
					Pax:        detailAdditional.Pax,
					IsRequired: detailAdditional.IsRequired,
				}

				additionals = append(additionals, cartAdditional)

				// Calculate total - only count price-based additionals
				if cartAdditional.Category == constant.AdditionalServiceCategoryPrice && cartAdditional.Price != nil {
					totalAdditional += *cartAdditional.Price
				}
			}

			var cancellationDate string
			cancellationDate = detail.CheckInDate.AddDate(0, 0, detail.RoomPrice.RoomType.Hotel.CancellationPeriod).Format("2006-01-02")

			nights := int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24)
			var checkInHourDur, checkOutHourDur time.Duration
			checkInHour := detail.RoomPrice.RoomType.Hotel.CheckInHour
			if checkInHour != nil {
				h, m, sec := checkInHour.Clock()
				checkInHourDur = time.Duration(h)*time.Hour +
					time.Duration(m)*time.Minute +
					time.Duration(sec)*time.Second
			}

			checkOutHour := detail.RoomPrice.RoomType.Hotel.CheckOutHour
			if checkOutHour != nil {
				h, m, sec := checkOutHour.Clock()
				checkOutHourDur = time.Duration(h)*time.Hour +
					time.Duration(m)*time.Minute +
					time.Duration(sec)*time.Second
			}

			cartDetail := bookingdto.CartDetail{
				ID:                   detail.ID,
				HotelName:            detail.RoomPrice.RoomType.Hotel.Name,
				HotelRating:          detail.RoomPrice.RoomType.Hotel.Rating,
				CheckInDate:          detail.CheckInDate.In(constant.AsiaJakarta).Add(checkInHourDur).Format(time.RFC3339),
				CheckOutDate:         detail.CheckOutDate.In(constant.AsiaJakarta).Add(checkOutHourDur).Format(time.RFC3339),
				RoomTypeName:         detail.RoomPrice.RoomType.Name,
				IsBreakfast:          detail.RoomPrice.IsBreakfast,
				Guest:                detail.Guest,
				BedTypes:             detail.BedTypeNames, // Available bed types for selection
				Additional:           additionals,
				CancellationDate:     cancellationDate,
				PriceBeforePromo:     detail.RoomPrice.Price * float64(nights),
				TotalAdditionalPrice: totalAdditional,
			}
			basePrice := detail.RoomPrice.Price
			roomPrice := basePrice
			if detail.Promo != nil {
				promo := detail.Promo
				detailPromo, err := bu.generateDetailPromo(promo)
				if err != nil {
					logger.Error(ctx, "failed to generate detail promo", err.Error())
				}
				cartDetail.Promo = detailPromo
				switch detail.Promo.PromoTypeID {
				case constant.PromoTypeFixedPriceID:
					roomPrice = promo.Detail.FixedPrice
					if promo.Duration > nights {
						roomPrice += float64(nights-promo.Duration) * basePrice
					}
				case constant.PromoTypeDiscountID:
					roomPrice = (100 - promo.Detail.DiscountPercentage) / 100 * basePrice * float64(nights)
					if promo.Duration > nights {
						roomPrice += float64(nights-promo.Duration) * basePrice
					}
				default:
					roomPrice = basePrice * float64(nights)
				}
			}
			cartDetail.Price = roomPrice
			cartDetail.TotalPrice = cartDetail.Price + cartDetail.TotalAdditionalPrice
			for _, photo := range detail.RoomPrice.RoomType.Photos {
				if photo != "" {
					bucketName := fmt.Sprintf("%s-%s", constant.ConstHotel, constant.ConstPublic)
					photoUrl, err := bu.fileStorage.GetFile(ctx, bucketName, photo)
					if err != nil {
						logger.Error(ctx, "ListHotelsForAgent", err.Error())
						continue
					}
					cartDetail.Photo = photoUrl
					break
				}
			}

			grandTotal += cartDetail.TotalPrice
			details = append(details, cartDetail)
		}

		result.Detail = details

		result.Guest = cart.Guests
		result.GrandTotal = grandTotal
	}

	return result, nil
}

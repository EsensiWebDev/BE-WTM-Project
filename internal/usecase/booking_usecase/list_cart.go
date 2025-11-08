package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
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
	grandTotal := int64(0)
	if cart != nil {
		details := make([]bookingdto.CartDetail, len(cart.BookingDetails))
		for i, detail := range cart.BookingDetails {

			additionals := make([]bookingdto.CartDetailAdditional, len(detail.BookingDetailAdditional))
			totalAdditional := 0
			for j, additional := range detail.BookingDetailAdditional {
				additionals[j] = bookingdto.CartDetailAdditional{
					Name:  additional.NameAdditional,
					Price: additional.Price,
				}
				totalAdditional += int(additional.Price)
			}

			details[i] = bookingdto.CartDetail{
				HotelName:    detail.DetailRooms.HotelName,
				HotelRating:  detail.DetailRooms.HotelRating,
				CheckInDate:  detail.CheckInDate,
				CheckOutDate: detail.CheckOutDate,
				RoomTypeName: detail.DetailRooms.RoomTypeName,
				IsBreakfast:  detail.DetailRooms.IsBreakfast,
				Guest:        detail.Guest,
				Additional:   additionals,
				Promo: bookingdto.CartDetailPromo{
					Type:            detail.DetailPromos.Type,
					DiscountPercent: detail.DetailPromos.DiscountPercent,
					FixedPrice:      detail.DetailPromos.FixedPrice,
					UpgradedToID:    detail.DetailPromos.UpgradedToID,
					Benefit:         detail.DetailPromos.BenefitNote,
				},
				Price:                detail.Price,
				TotalAdditionalPrice: float64(totalAdditional),
				TotalPrice:           detail.Price + float64(totalAdditional),
			}
			grandTotal += int64(details[i].TotalPrice)
		}
		result.Detail = details
		result.Guest = cart.Guests
		result.GrandTotal = grandTotal
	}

	return result, nil
}

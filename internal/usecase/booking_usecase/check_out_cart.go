package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) CheckOutCart(ctx context.Context, req *bookingdto.CheckOutCartRequest) error {
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

		booking, err := bu.bookingRepo.GetCartBooking(txCtx, agentID)
		if err != nil {
			logger.Error(ctx, "failed to get card booking", err.Error())
			return fmt.Errorf("failed to get cart: %s", err.Error())
		}
		if booking == nil {
			logger.Error(ctx, "card booking is nil")
			return fmt.Errorf("no cart found")
		}

		// 1. Update guests in booking
		if err := bu.bookingRepo.UpdateBookingGuests(txCtx, booking.ID, req.Guests); err != nil {
			logger.Error(ctx, "failed to update booking", err.Error())
			return fmt.Errorf("failed to update booking guests: %s", err.Error())
		}

		// 2. Update guest per booking detail
		for _, d := range req.Details {
			if err := bu.bookingRepo.UpdateBookingDetailGuest(txCtx, d.BookingDetailID, d.Guest); err != nil {
				logger.Error(ctx, "failed to update booking detail", err.Error())
				return fmt.Errorf("failed to update guest in detail: %s", err.Error())
			}
		}

		// 3. Update status to "in review"
		if err := bu.bookingRepo.UpdateBookingStatus(txCtx, booking.ID, constant.StatusBookingInReviewID); err != nil {
			logger.Error(ctx, "failed to update booking status", err.Error())
			return fmt.Errorf("failed to update booking status: %s", err.Error())
		}

		return nil
	})
}

package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) RemoveFromCart(ctx context.Context, bookingDetailID uint) error {
	// Get agent Id from context
	userCtx, err := bu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return fmt.Errorf("failed to get user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return fmt.Errorf("user context is nil")
	}

	agentID := userCtx.ID

	// Remove BookingDetail from cart
	if err := bu.bookingRepo.DeleteCartBooking(ctx, agentID, bookingDetailID); err != nil {
		logger.Error(ctx, "failed to remove booking detail from cart", err.Error())
		return fmt.Errorf("failed to remove booking detail from cart: %s", err.Error())
	}

	return nil
}

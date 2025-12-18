package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) RemoveGuestsFromCart(ctx context.Context, req *bookingdto.RemoveGuestsFromCartRequest) error {
	// Get agent Id from context
	userCtx, err := bu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return fmt.Errorf("failed to get user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return fmt.Errorf("user not found in context")
	}

	agentID := userCtx.ID
	if err := bu.bookingRepo.RemoveGuestsFromCart(ctx, agentID, req.CartID, req.Guests); err != nil {
		logger.Error(ctx, "failed to remove guests from cart", err.Error())
		return err
	}

	return nil
}

package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) AddGuestsToCart(ctx context.Context, req *bookingdto.AddGuestsToCartRequest) error {

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

	if err := bu.bookingRepo.AddGuestsToCart(ctx, agentID, req.CartID, req.Guests); err != nil {
		logger.Error(ctx, "failed to add guests to cart", err.Error())
		return err
	}

	return nil
}

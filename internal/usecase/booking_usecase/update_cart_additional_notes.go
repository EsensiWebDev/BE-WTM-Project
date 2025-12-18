package booking_usecase

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/logger"
)

// UpdateCartAdditionalNotes updates additional_notes for a specific sub-cart item (booking_detail)
// that belongs to the authenticated agent's current cart.
func (bu *BookingUsecase) UpdateCartAdditionalNotes(ctx context.Context, req *bookingdto.UpdateCartAdditionalNotesRequest) error {
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

	// Trim whitespace on backend side as an extra safety layer
	notes := strings.TrimSpace(req.AdditionalNotes)

	if err := bu.bookingRepo.UpdateCartAdditionalNotes(ctx, agentID, req.SubCartID, notes); err != nil {
		logger.Error(ctx, "failed to update cart additional notes", err.Error())
		return fmt.Errorf("failed to update cart additional notes: %s", err.Error())
	}

	return nil
}




package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListSubBookingIDs(ctx context.Context, req *bookingdto.ListSubBookingIDsRequest) (*bookingdto.ListSubBookingIDsResponse, error) {
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

	subBookingIDs, err := bu.bookingRepo.GetSubBookingIDs(ctx, agentID, req.BookingID)
	if err != nil {
		logger.Error(ctx, "failed to get sub booking IDs by agent ID", err.Error())
		return nil, fmt.Errorf("failed to get sub booking IDs by agent ID: %s", err.Error())
	}

	return &bookingdto.ListSubBookingIDsResponse{
		SubBookingIDs: subBookingIDs,
	}, nil
}

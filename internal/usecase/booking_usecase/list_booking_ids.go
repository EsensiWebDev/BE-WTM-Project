package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListBookingIDs(ctx context.Context, req *bookingdto.ListBookingIDsRequest) (*bookingdto.ListBookingIDsResponse, error) {
	filterReq := filter.DefaultFilter{}
	filterReq.PaginationRequest = req.PaginationRequest

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

	bookingIDs, total, err := bu.bookingRepo.GetBookingIDs(ctx, agentID, &filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get booking IDs", err.Error())
		return nil, err
	}

	resp := &bookingdto.ListBookingIDsResponse{
		Total:      total,
		BookingIDs: bookingIDs,
	}

	return resp, nil
}

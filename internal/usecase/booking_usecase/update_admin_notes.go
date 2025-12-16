package booking_usecase

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/logger"
)

// UpdateAdminNotes updates admin_notes for a booking detail identified by sub_booking_id.
// This is used by admins to add notes for agents to read.
func (bu *BookingUsecase) UpdateAdminNotes(ctx context.Context, req *bookingdto.UpdateAdminNotesRequest) error {
	// Trim whitespace on backend side as an extra safety layer
	notes := strings.TrimSpace(req.AdminNotes)

	if err := bu.bookingRepo.UpdateAdminNotes(ctx, req.SubBookingID, notes); err != nil {
		logger.Error(ctx, "failed to update admin notes", err.Error())
		return fmt.Errorf("failed to update admin notes: %s", err.Error())
	}

	return nil
}

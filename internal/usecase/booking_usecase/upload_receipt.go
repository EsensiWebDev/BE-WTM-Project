package booking_usecase

import (
	"context"
	"strings"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) UploadReceipt(ctx context.Context, req *bookingdto.UploadReceiptRequest) error {
	var bookindDetailIDs []uint
	var prefix string
	var id uint

	if strings.TrimSpace(req.BookingID) != "" {
		prefix = "booking/receipts/booking"

		booking, err := bu.bookingRepo.GetBookingByCode(ctx, req.BookingID)
		if err != nil {
			logger.Error(ctx, "failed to get bookings", err.Error())
			return err
		}
		id = booking.ID
		bookindDetailIDs = make([]uint, 0, len(booking.BookingDetails))
		for _, detail := range booking.BookingDetails {
			bookindDetailIDs = append(bookindDetailIDs, detail.ID)
		}
	} else {
		detail, err := bu.bookingRepo.GetSubBookingByCode(ctx, req.BookingDetailID)
		if err != nil {
			logger.Error(ctx, "failed to get sub booking by code", err.Error())
			return err
		}
		bookindDetailIDs = append(bookindDetailIDs, detail.ID)
		prefix = "booking/receipts/booking_detail"
		id = detail.BookingID
	}

	fileReceiptPath, err := bu.uploadFile(ctx, req.Receipt, prefix, id)
	if err != nil {
		logger.Error(ctx, "failed to upload receipt file", err.Error())
		return err
	}

	if err := bu.bookingRepo.UpdateBookingReceipt(ctx, bookindDetailIDs, fileReceiptPath); err != nil {
		logger.Error(ctx, "failed to update booking receipt", err.Error())
		return err
	}

	return nil
}

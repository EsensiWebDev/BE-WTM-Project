package booking_usecase

import (
	"context"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) UploadReceipt(ctx context.Context, req *bookingdto.UploadReceiptRequest) error {
	var bookindDetailIDs []uint
	var prefix string
	var id uint

	if req.BookingID > 0 {
		prefix = "booking/receipts/booking"
		id = req.BookingID

		booking, err := bu.bookingRepo.GetBookingByID(ctx, req.BookingID)
		if err != nil {
			logger.Error(ctx, "failed to get bookings", err.Error())
			return err
		}
		bookindDetailIDs = make([]uint, 0, len(booking.BookingDetails))
		for _, detail := range booking.BookingDetails {
			bookindDetailIDs = append(bookindDetailIDs, detail.ID)
		}
	} else {
		bookindDetailIDs = append(bookindDetailIDs, req.BookingDetailID)
		prefix = "booking/receipts/booking_detail"
		id = req.BookingDetailID
	}

	fileReceiptPath, err := bu.uploadFile(ctx, req.FileReceipt, prefix, id)
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

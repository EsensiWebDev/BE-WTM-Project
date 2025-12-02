package booking_repository

import (
	"context"
	"errors"
)

func (br *BookingRepository) GetBookingDetailIDsByBookingCode(ctx context.Context, bookingCode string) ([]uint, error) {
	if bookingCode == "" {
		return nil, errors.New("booking code cannot be empty")
	}

	db := br.db.GetTx(ctx)
	query := `SELECT bd.id 
			FROM booking_details bd
			JOIN bookings b ON bd.booking_id = b.id
			WHERE b.booking_code = ?`

	var ids []uint
	if err := db.WithContext(ctx).
		Raw(query, bookingCode).
		Scan(&ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

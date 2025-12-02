package booking_repository

import (
	"context"
	"errors"
)

func (br *BookingRepository) GetIDBySubBookingID(ctx context.Context, subBookingID string) (uint, error) {
	if subBookingID == "" {
		return 0, errors.New("sub booking id cannot be empty")
	}

	db := br.db.GetTx(ctx)
	query := `SELECT id FROM booking_details WHERE sub_booking_id = ?`

	var id uint
	if err := db.WithContext(ctx).
		Raw(query, subBookingID).
		Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

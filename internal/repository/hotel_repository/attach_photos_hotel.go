package hotel_repository

import (
	"context"
	"github.com/lib/pq"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) AttachPhotosHotel(ctx context.Context, hotelID uint, photoURLs []string) error {
	db := hr.db.GetTx(ctx)

	if err := db.WithContext(ctx).Model(&model.Hotel{}).Where("id = ?", hotelID).
		Updates(map[string]interface{}{
			"photos": pq.StringArray(photoURLs),
		}).Error; err != nil {
		logger.Error(ctx, "Failed to attach photos to hotel", err.Error())
		return err
	}

	return nil
}

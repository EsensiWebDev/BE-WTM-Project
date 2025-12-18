package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) GetRoomTypePreferencesByIDs(ctx context.Context, ids []uint) ([]entity.RoomTypePreference, error) {
	db := hr.db.GetTx(ctx)

	if len(ids) == 0 {
		return nil, nil
	}
	var preferences []model.RoomTypePreference
	if err := db.WithContext(ctx).
		Preload("OtherPreference").
		Where("id IN ?", ids).
		Find(&preferences).Error; err != nil {
		if hr.db.ErrRecordNotFound(ctx, err) {
			logger.Error(ctx, "Not found")
			return nil, nil
		}
		return nil, err
	}

	// Convert model.RoomTypePreference to entity.RoomTypePreference
	var preferencesEntity []entity.RoomTypePreference
	for _, pref := range preferences {
		preferencesEntity = append(preferencesEntity, entity.RoomTypePreference{
			ID:                pref.ID,
			RoomTypeID:        pref.RoomTypeID,
			OtherPreferenceID: pref.OtherPreferenceID,
			OtherPreference: entity.OtherPreference{
				ID:   pref.OtherPreference.ID,
				Name: pref.OtherPreference.Name,
			},
		})
	}

	return preferencesEntity, nil
}

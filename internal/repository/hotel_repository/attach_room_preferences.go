package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

// AttachRoomPreferences attaches a list of "Other Preferences" (by name) to a room type.
// It will first-or-create the master OtherPreference by name, then create the link in RoomTypePreference.
func (hr *HotelRepository) AttachRoomPreferences(ctx context.Context, roomTypeID uint, preferenceNames []string) error {
	db := hr.db.GetTx(ctx)

	for _, rawName := range preferenceNames {
		name := strings.TrimSpace(rawName)
		if name == "" {
			continue
		}

		op := model.OtherPreference{Name: name}
		if err := db.Where("name = ?", name).FirstOrCreate(&op).Error; err != nil {
			logger.Error(ctx, "Failed to create other preference", err.Error())
			return err
		}

		link := model.RoomTypePreference{
			RoomTypeID:        roomTypeID,
			OtherPreferenceID: op.ID,
		}
		if err := db.Create(&link).Error; err != nil {
			logger.Error(ctx, "Failed to attach other preference to room type", err.Error())
			return err
		}
	}

	return nil
}

// UpdateRoomPreferences updates the "Other Preferences" attached to a room type.
// - It keeps links whose IDs are in unchangedPreferenceIDs.
// - It deletes other existing links for this room type.
// - It attaches new preferences from newPreferenceNames.
// - It also cleans up orphan OtherPreference records that are no longer referenced.
func (hr *HotelRepository) UpdateRoomPreferences(ctx context.Context, roomTypeID uint, unchangedPreferenceIDs []uint, newPreferenceNames []string) error {
	db := hr.db.GetTx(ctx)

	// Delete RoomTypePreference links that are no longer attached to this room type
	deleteQuery := db.WithContext(ctx).
		Unscoped().
		Where("room_type_id = ?", roomTypeID)

	if len(unchangedPreferenceIDs) > 0 {
		deleteQuery = deleteQuery.Where("id NOT IN (?)", unchangedPreferenceIDs)
	}

	if err := deleteQuery.Delete(&model.RoomTypePreference{}).Error; err != nil {
		logger.Error(ctx, "Failed to delete existing room type preferences", err.Error())
		return err
	}

	// Attach new preferences
	if err := hr.AttachRoomPreferences(ctx, roomTypeID, newPreferenceNames); err != nil {
		return err
	}

	// Clean up orphan OtherPreference records that are no longer referenced
	var orphanPreferenceIDs []uint
	if err := db.WithContext(ctx).
		Model(&model.OtherPreference{}).
		Joins("LEFT JOIN room_type_preferences rtp ON rtp.other_preference_id = other_preferences.id").
		Where("rtp.id IS NULL").
		Pluck("other_preferences.id", &orphanPreferenceIDs).Error; err != nil {
		logger.Error(ctx, "Failed to find orphan other preferences", err.Error())
		return err
	}

	if len(orphanPreferenceIDs) > 0 {
		if err := db.WithContext(ctx).
			Unscoped().
			Where("id IN (?)", orphanPreferenceIDs).
			Delete(&model.OtherPreference{}).Error; err != nil {
			logger.Error(ctx, "Failed to delete orphan other preferences", err.Error())
			return err
		}
	}

	return nil
}

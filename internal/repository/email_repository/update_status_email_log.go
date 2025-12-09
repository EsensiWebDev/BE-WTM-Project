package email_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"

	"gorm.io/gorm"
)

func (er *EmailRepository) UpdateStatusEmailLog(ctx context.Context, log *entity.EmailLog) error {
	db := er.db.GetTx(ctx)

	id := log.ID
	statusID := log.StatusID
	meta := log.Meta

	// Build update fields
	updateFields := map[string]interface{}{
		"status_id":  statusID,
		"updated_at": gorm.Expr("NOW()"),
	}

	// Jika meta tidak nil, tambahkan ke update
	if meta != nil {
		// Marshal meta ke JSON jika perlu (tergantung bagaimana Anda menyimpan di DB)
		metaJSON, err := json.Marshal(meta)
		if err != nil {
			logger.Error(ctx, "Error marshaling email log meta", err.Error())
			return fmt.Errorf("failed to marshal meta: %w", err)
		}
		updateFields["meta"] = metaJSON
	}

	// Execute update
	result := db.Model(&model.EmailLog{}).
		Where("id = ?", id).
		Updates(updateFields)

	if result.Error != nil {
		logger.Error(ctx, "Error updating status email log", result.Error.Error())
		return fmt.Errorf("failed to update email log: %w", result.Error)
	}

	// Cek jika record tidak ditemukan
	if result.RowsAffected == 0 {
		logger.Warn(ctx, "Email log not found with id: %d", id)
		return fmt.Errorf("email log with id %d not found", id)
	}

	logger.Info(ctx, "Successfully updated email log id: %d, status: %d", id, statusID)
	return nil
}

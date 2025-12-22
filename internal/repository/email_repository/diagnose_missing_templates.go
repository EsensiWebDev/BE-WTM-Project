package email_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

// DiagnoseMissingTemplates finds email logs with missing email templates
// This can be called to identify orphaned email logs
func (er *EmailRepository) DiagnoseMissingTemplates(ctx context.Context) error {
	db := er.db.GetTx(ctx)

	// Find email logs where the email_template_id doesn't exist in email_templates
	var orphanedLogs []struct {
		ID              uint
		EmailTemplateID uint
		CreatedAt       string
	}

	query := `
		SELECT el.id, el.email_template_id, el.created_at::text
		FROM email_logs el
		LEFT JOIN email_templates et ON el.email_template_id = et.id
		WHERE et.id IS NULL
		ORDER BY el.created_at DESC
		LIMIT 100
	`

	if err := db.WithContext(ctx).Raw(query).Scan(&orphanedLogs).Error; err != nil {
		logger.Error(ctx, "Failed to diagnose missing templates", err.Error())
		return err
	}

	if len(orphanedLogs) > 0 {
		logger.Error(ctx, "Found %d email logs with missing templates:", len(orphanedLogs))
		for _, log := range orphanedLogs {
			logger.Error(ctx, "  - Email Log ID: %d references missing Template ID: %d (created at: %s)", log.ID, log.EmailTemplateID, log.CreatedAt)
		}
	} else {
		logger.Info(ctx, "No orphaned email logs found - all templates exist")
	}

	// Also check what templates exist
	var existingTemplates []model.EmailTemplate
	if err := db.WithContext(ctx).Find(&existingTemplates).Error; err != nil {
		logger.Error(ctx, "Failed to list existing templates", err.Error())
	} else {
		logger.Info(ctx, "Existing email templates in database:")
		for _, tpl := range existingTemplates {
			logger.Info(ctx, "  - Template ID: %d, Name: %s", tpl.ID, tpl.Name)
		}
	}

	return nil
}

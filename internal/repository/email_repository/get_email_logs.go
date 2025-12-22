package email_repository

import (
	"context"
	"encoding/json"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"

	"gorm.io/gorm/clause"
)

func (er *EmailRepository) GetEmailLogs(ctx context.Context, filter filter.EmailLogFilter) ([]entity.EmailLog, int64, error) {
	db := er.db.GetTx(ctx)
	query := db.WithContext(ctx).
		Model(&model.EmailLog{})

	// Apply email type filter using subquery instead of JOIN to avoid Preload conflicts
	if len(filter.EmailType) > 0 {
		query = query.Where("email_logs.email_template_id IN (SELECT id FROM email_templates WHERE name IN ?)", filter.EmailType)
	}

	// Apply status filter
	if len(filter.Status) > 0 {
		statusIDs := make([]uint, 0)
		for _, status := range filter.Status {
			switch status {
			case constant.StatusEmailPending:
				statusIDs = append(statusIDs, constant.StatusEmailPendingID)
			case constant.StatusEmailSuccess:
				statusIDs = append(statusIDs, constant.StatusEmailSuccessID)
			case constant.StatusEmailFailed:
				statusIDs = append(statusIDs, constant.StatusEmailFailedID)
			}
		}
		if len(statusIDs) > 0 {
			query = query.Where("email_logs.status_id IN ?", statusIDs)
		}
	}

	// Apply hotel name filter (search in meta JSON)
	if filter.HotelName != "" {
		query = query.Where("email_logs.meta->>'hotel_name' ILIKE ?", "%"+filter.HotelName+"%")
	}

	// Apply booking code filter (search in meta JSON)
	if filter.BookingCode != "" {
		query = query.Where("email_logs.meta->>'booking_code' ILIKE ?", "%"+filter.BookingCode+"%")
	}

	// Apply date range filter
	if filter.DateFrom != nil && !filter.DateFrom.IsZero() {
		query = query.Where("email_logs.created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != nil && !filter.DateTo.IsZero() {
		query = query.Where("email_logs.created_at <= ?", filter.DateTo)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "failed to count email logs", err.Error())
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	// Apply sorting
	// Map frontend column names to database column names
	columnMap := map[string]string{
		"date_time":  "created_at",
		"created_at": "created_at",
	}

	// Default sort if no sort specified
	defaultSort := "created_at"
	defaultDir := "desc"

	if filter.Sort != "" {
		// Check if the sort column is allowed and map it
		if dbColumn, exists := columnMap[filter.Sort]; exists {
			dir := strings.TrimSpace(strings.ToLower(filter.Dir))
			desc := dir != "asc"
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: dbColumn}, Desc: desc})
		} else {
			// Invalid sort column, use default
			query = query.Order(defaultSort + " " + defaultDir)
		}
	} else {
		// No sort specified, use default
		query = query.Order(defaultSort + " " + defaultDir)
	}

	// Fetch results
	// Use Preload with error handling - if Preload fails, templates will have ID = 0
	emailLogs := []model.EmailLog{}
	if err := query.Preload("EmailTemplate").Find(&emailLogs).Error; err != nil {
		logger.Error(ctx, "failed to get email logs", err.Error())
		return nil, 0, err
	}

	// Diagnostic: Check for missing templates and log them
	// Also verify Preload worked correctly and manually fix if needed
	for i, emailLog := range emailLogs {
		if emailLog.EmailTemplate.ID == 0 {
			logger.Error(ctx, "❌ MISSING TEMPLATE - Email Log ID: %d references Template ID: %d which does not exist or Preload failed", emailLog.ID, emailLog.EmailTemplateID)
			// Try to manually verify if template exists (Preload might have failed even if template exists)
			var templateExists model.EmailTemplate
			if err := db.WithContext(ctx).Where("id = ?", emailLog.EmailTemplateID).First(&templateExists).Error; err != nil {
				logger.Error(ctx, "  → Confirmed: Template ID %d does NOT exist in database", emailLog.EmailTemplateID)
			} else {
				logger.Error(ctx, "  → Template ID %d EXISTS but Preload failed! Template name: %s", emailLog.EmailTemplateID, templateExists.Name)
				// Manually set the template since Preload failed
				emailLogs[i].EmailTemplate = templateExists
			}
		}
	}

	// Convert to entity.EmailLog slice
	// Use defensive copying to avoid panics from nil templates
	var result []entity.EmailLog
	result = make([]entity.EmailLog, len(emailLogs))

	// Manually copy each field to avoid issues with CopyStrict and nil templates
	for i, emailLog := range emailLogs {
		result[i] = entity.EmailLog{
			ID:              emailLog.ID,
			To:              emailLog.To,
			Subject:         emailLog.Subject,
			Body:            emailLog.Body,
			EmailTemplateID: emailLog.EmailTemplateID,
			StatusID:        emailLog.StatusID,
			CreatedAt:       emailLog.CreatedAt,
		}

		// Copy meta if exists
		if emailLog.Meta != nil {
			var meta entity.MetadataEmailLog
			if err := json.Unmarshal(emailLog.Meta, &meta); err != nil {
				logger.Error(ctx, "failed to unmarshal email log meta", err.Error())
				return nil, 0, err
			}
			result[i].Meta = &meta
		}
	}

	for i, emailLog := range emailLogs {
		if emailLog.Meta != nil {
			var meta entity.MetadataEmailLog
			if err := json.Unmarshal(emailLog.Meta, &meta); err != nil {
				logger.Error(ctx, "failed to unmarshal email log meta", err.Error())
				return nil, 0, err
			}
			result[i].Meta = &meta
		}

		// Safely get email type from template
		// Use a helper to safely access template name to avoid nil pointer dereference
		result[i].EmailType = getEmailTypeFromTemplate(ctx, emailLog)
	}

	return result, total, nil
}

// getEmailTypeFromTemplate safely extracts email type from EmailTemplate
// Returns empty string if template is not found or invalid
// Uses recover to handle any potential nil pointer panics
func getEmailTypeFromTemplate(ctx context.Context, emailLog model.EmailLog) string {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(ctx, "Panic recovered in getEmailTypeFromTemplate for email log ID: %d, template ID: %d, error: %v", emailLog.ID, emailLog.EmailTemplateID, r)
		}
	}()

	// Check if EmailTemplate was successfully loaded by GORM Preload
	// When Preload fails, the struct will be zero-valued (ID = 0)
	if emailLog.EmailTemplate.ID == 0 {
		logger.Error(ctx, "❌ MISSING EMAIL TEMPLATE - Email log ID: %d references template ID: %d which does not exist in database", emailLog.ID, emailLog.EmailTemplateID)
		return ""
	}

	// Safely get template name - access Name field carefully
	var templateName string
	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, "Panic accessing EmailTemplate.Name for email log ID: %d, template ID: %d, error: %v", emailLog.ID, emailLog.EmailTemplateID, r)
				templateName = ""
			}
		}()
		templateName = emailLog.EmailTemplate.Name
	}()

	if templateName == "" {
		logger.Error(ctx, "❌ EMPTY TEMPLATE NAME - Email log ID: %d has template ID: %d but template name is empty", emailLog.ID, emailLog.EmailTemplateID)
		return ""
	}

	// Look up email type in mapping
	if emailType, exists := constant.MapEmailType[templateName]; exists {
		return emailType
	}

	logger.Warn(ctx, "Email type mapping not found for template name: '%s' (email log ID: %d, template ID: %d). Available mappings: %v", templateName, emailLog.ID, emailLog.EmailTemplateID, getMapKeys(constant.MapEmailType))
	return ""
}

// getMapKeys returns all keys from a map (helper for logging)
func getMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

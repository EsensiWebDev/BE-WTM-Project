// database/migration.go
package database

import (
	"context"
	"fmt"
	"reflect"
	"wtm-backend/config"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/infrastructure/database/seed"
	"wtm-backend/pkg/logger"
)

func (dbs *DBPostgre) runMigrations(ctx context.Context, cfg *config.Config) error {
	// Skip auto-migrate di production kecuali explicit di enable
	if cfg.IsProduction() && !cfg.AutoMigrate {
		logger.Info(ctx, "Auto-migration skipped in production")
		return nil
	}

	logger.Info(ctx, "Running database migrations")

	models := []interface{}{
		&model.StatusHotel{},
		&model.Hotel{},
		&model.NearbyPlace{},
		&model.HotelNearbyPlace{},
		&model.Facility{},
		&model.RoomType{},
		&model.BedType{},
		&model.RoomAdditional{},
		&model.RoomTypeAdditional{},
		&model.RoomUnavailable{},
		&model.PromoType{},
		&model.Promo{},
		&model.PromoGroup{},
		&model.PromoRoomType{},
		&model.User{},
		&model.StatusUser{},
		&model.AgentCompany{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
		&model.RoomPrice{},
		&model.Banner{},
		&model.StatusBooking{},
		&model.StatusPayment{},
		&model.Booking{},
		&model.BookingDetail{},
		&model.BookingDetailAdditional{},
		&model.BookingGuest{},
		&model.EmailTemplate{},
		&model.Notification{},
		&model.UserNotificationSetting{},
		&model.PasswordResetToken{},
		&model.StatusEmail{},
		&model.EmailLog{},
		&model.Invoice{},
	}

	if err := dbs.DB.AutoMigrate(models...); err != nil {
		logger.Error(ctx, "Database migration failed", err.Error())
		return fmt.Errorf("migration: %w", err)
	}

	// ✅ TAMBAHKAN: Migrasi external_id untuk semua table
	if err := dbs.migrateAllExternalIDs(ctx, models); err != nil {
		logger.Error(ctx, "ExternalID migration failed", err.Error())
		return fmt.Errorf("external_id migration: %w", err)
	}

	logger.Info(ctx, "Database migration completed",
		fmt.Sprintf("models: %d", len(models)))

	// Seeding hanya di non-production
	if !cfg.IsProduction() {
		dbs.runSeeding(ctx)
	}

	return nil
}

func (dbs *DBPostgre) runSeeding(ctx context.Context) {
	if err := seed.Seeding(dbs.DB); err != nil {
		logger.Warn(ctx, "Database seeding failed", err.Error())
	} else {
		logger.Info(ctx, "Database seeding completed")
	}
}

// ✅ FUNGSI BARU: Migrasi external_id untuk semua model
func (dbs *DBPostgre) migrateAllExternalIDs(ctx context.Context, models []interface{}) error {
	logger.Info(ctx, "Starting ExternalID migration for all tables")

	for _, model := range models {
		tableName := dbs.DB.NamingStrategy.TableName(reflect.TypeOf(model).Elem().Name())

		if err := dbs.migrateTableExternalID(ctx, tableName, model); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", tableName, err)
		}
	}

	logger.Info(ctx, "ExternalID migration completed successfully")
	return nil
}

// ✅ FUNGSI BARU: Migrasi external_id per table (Pure SQL)
func (dbs *DBPostgre) migrateTableExternalID(ctx context.Context, tableName string, model interface{}) error {
	if tableName == "role_permissions" {
		// Skip role_permissions migration
		return nil
	}

	// Step 1: Check if external_id column already exists menggunakan SQL
	var columnExists bool
	checkColumnSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'external_id'
		)
	`
	if err := dbs.DB.Raw(checkColumnSQL, tableName).Scan(&columnExists).Error; err != nil {
		return fmt.Errorf("failed to check column existence in %s: %w", tableName, err)
	}

	if columnExists {
		logger.Info(ctx, fmt.Sprintf("ExternalID column already exists in %s", tableName))
	} else {
		// Step 2: Add external_id column (nullable first)
		logger.Info(ctx, fmt.Sprintf("Adding ExternalID column to %s", tableName))

		addColumnSQL := fmt.Sprintf(`
			ALTER TABLE %s ADD COLUMN external_id TEXT
		`, tableName)

		if err := dbs.DB.Exec(addColumnSQL).Error; err != nil {
			return fmt.Errorf("failed to add column to %s: %w", tableName, err)
		}
	}

	// Step 3: Backfill existing data dengan pure UUID
	logger.Info(ctx, fmt.Sprintf("Backfilling ExternalID for existing data in %s", tableName))

	backfillSQL := fmt.Sprintf(`
		UPDATE %s 
		SET external_id = gen_random_uuid()::text
		WHERE external_id IS NULL OR external_id = ''
	`, tableName)

	if err := dbs.DB.Exec(backfillSQL).Error; err != nil {
		return fmt.Errorf("failed to backfill %s: %w", tableName, err)
	}

	// Step 4: Add NOT NULL constraint setelah data terisi
	logger.Info(ctx, fmt.Sprintf("Adding NOT NULL constraint to %s", tableName))

	constraintSQL := fmt.Sprintf(`
		ALTER TABLE %s 
		ALTER COLUMN external_id SET NOT NULL
	`, tableName)

	if err := dbs.DB.Exec(constraintSQL).Error; err != nil {
		logger.Warn(ctx, fmt.Sprintf("NOT NULL constraint might already exist on %s: %v", tableName, err))
	}

	// Step 5: Verify no NULL values
	var nullCount int64
	verifySQL := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE external_id IS NULL OR external_id = ''`, tableName)
	if err := dbs.DB.Raw(verifySQL).Scan(&nullCount).Error; err != nil {
		return fmt.Errorf("failed to verify %s: %w", tableName, err)
	}

	if nullCount > 0 {
		return fmt.Errorf("%s still has %d NULL external_id values", tableName, nullCount)
	}

	logger.Info(ctx, fmt.Sprintf("✓ Successfully migrated ExternalID for %s", tableName))
	return nil
}

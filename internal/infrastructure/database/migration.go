// database/migration.go
package database

import (
	"context"
	"fmt"
	"reflect"
	"wtm-backend/config"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/infrastructure/database/seed"
	"wtm-backend/pkg/constant"
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
		&model.OtherPreference{},
		&model.RoomTypePreference{},
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
		&model.Currency{},
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

	// ✅ Migrate RoomTypeAdditional new fields (Category, Pax, IsRequired)
	if err := dbs.migrateRoomTypeAdditional(ctx); err != nil {
		logger.Error(ctx, "RoomTypeAdditional migration failed", err.Error())
		return fmt.Errorf("room_type_additional migration: %w", err)
	}

	// ✅ Migrate BookingDetailAdditional new fields (Category, Pax, IsRequired, nullable Price)
	if err := dbs.migrateBookingDetailAdditional(ctx); err != nil {
		logger.Error(ctx, "BookingDetailAdditional migration failed", err.Error())
		return fmt.Errorf("booking_detail_additional migration: %w", err)
	}

	// ✅ Migrate BookingGuest new fields (Honorific, Category, Age)
	if err := dbs.migrateBookingGuest(ctx); err != nil {
		logger.Error(ctx, "BookingGuest migration failed", err.Error())
		return fmt.Errorf("booking_guest migration: %w", err)
	}

	// ✅ Migrate BookingDetail bed_type field
	if err := dbs.migrateBookingDetailBedType(ctx); err != nil {
		logger.Error(ctx, "BookingDetail bed_type migration failed", err.Error())
		return fmt.Errorf("booking_detail bed_type migration: %w", err)
	}

	// ✅ Migrate BookingDetail additional_notes field
	if err := dbs.migrateBookingDetailAdditionalNotes(ctx); err != nil {
		logger.Error(ctx, "BookingDetail additional_notes migration failed", err.Error())
		return fmt.Errorf("booking_detail additional_notes migration: %w", err)
	}

	// ✅ Migrate BookingDetail admin_notes field
	if err := dbs.migrateBookingDetailAdminNotes(ctx); err != nil {
		logger.Error(ctx, "BookingDetail admin_notes migration failed", err.Error())
		return fmt.Errorf("booking_detail admin_notes migration: %w", err)
	}

	// ✅ Migrate Multi-Currency support
	if err := dbs.migrateMultiCurrency(ctx); err != nil {
		logger.Error(ctx, "Multi-currency migration failed", err.Error())
		return fmt.Errorf("multi-currency migration: %w", err)
	}

	// ✅ Migrate Currency Symbol: Update IDR symbol from "Rp" to "IDR"
	if err := dbs.migrateCurrencySymbol(ctx); err != nil {
		logger.Error(ctx, "Currency symbol migration failed", err.Error())
		return fmt.Errorf("currency symbol migration: %w", err)
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

// ✅ FUNGSI BARU: Migrasi RoomTypeAdditional untuk menambahkan Category, Pax, IsRequired
func (dbs *DBPostgre) migrateRoomTypeAdditional(ctx context.Context) error {
	logger.Info(ctx, "Starting RoomTypeAdditional migration")

	tableName := "room_type_additionals"

	// Step 1: Check and add Category column
	var categoryExists bool
	checkCategorySQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'category'
		)
	`
	if err := dbs.DB.Raw(checkCategorySQL, tableName).Scan(&categoryExists).Error; err != nil {
		return fmt.Errorf("failed to check category column: %w", err)
	}

	if !categoryExists {
		logger.Info(ctx, "Adding category column to room_type_additionals")
		addCategorySQL := fmt.Sprintf(`
			ALTER TABLE room_type_additionals 
			ADD COLUMN category VARCHAR(10) DEFAULT '%s'
		`, "price")
		if err := dbs.DB.Exec(addCategorySQL).Error; err != nil {
			return fmt.Errorf("failed to add category column: %w", err)
		}
		// Set default value for existing records
		updateCategorySQL := fmt.Sprintf(`
			UPDATE room_type_additionals 
			SET category = '%s' 
			WHERE category IS NULL OR category = ''
		`, "price")
		if err := dbs.DB.Exec(updateCategorySQL).Error; err != nil {
			return fmt.Errorf("failed to set default category: %w", err)
		}
	}

	// Step 2: Check and add Pax column
	var paxExists bool
	checkPaxSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'pax'
		)
	`
	if err := dbs.DB.Raw(checkPaxSQL, tableName).Scan(&paxExists).Error; err != nil {
		return fmt.Errorf("failed to check pax column: %w", err)
	}

	if !paxExists {
		logger.Info(ctx, "Adding pax column to room_type_additionals")
		addPaxSQL := `
			ALTER TABLE room_type_additionals 
			ADD COLUMN pax INTEGER
		`
		if err := dbs.DB.Exec(addPaxSQL).Error; err != nil {
			return fmt.Errorf("failed to add pax column: %w", err)
		}
	}

	// Step 3: Check and add IsRequired column
	var isRequiredExists bool
	checkIsRequiredSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'is_required'
		)
	`
	if err := dbs.DB.Raw(checkIsRequiredSQL, tableName).Scan(&isRequiredExists).Error; err != nil {
		return fmt.Errorf("failed to check is_required column: %w", err)
	}

	if !isRequiredExists {
		logger.Info(ctx, "Adding is_required column to room_type_additionals")
		addIsRequiredSQL := `
			ALTER TABLE room_type_additionals 
			ADD COLUMN is_required BOOLEAN DEFAULT false
		`
		if err := dbs.DB.Exec(addIsRequiredSQL).Error; err != nil {
			return fmt.Errorf("failed to add is_required column: %w", err)
		}
		// Set default value for existing records
		updateIsRequiredSQL := `
			UPDATE room_type_additionals 
			SET is_required = false 
			WHERE is_required IS NULL
		`
		if err := dbs.DB.Exec(updateIsRequiredSQL).Error; err != nil {
			return fmt.Errorf("failed to set default is_required: %w", err)
		}
	}

	// Step 4: Make Price nullable if it's not already
	var priceNullable bool
	checkPriceNullableSQL := `
		SELECT is_nullable = 'YES'
		FROM information_schema.columns 
		WHERE table_name = $1 AND column_name = 'price'
	`
	if err := dbs.DB.Raw(checkPriceNullableSQL, tableName).Scan(&priceNullable).Error; err != nil {
		return fmt.Errorf("failed to check price nullable: %w", err)
	}

	if !priceNullable {
		logger.Info(ctx, "Making price column nullable in room_type_additionals")
		alterPriceSQL := `
			ALTER TABLE room_type_additionals 
			ALTER COLUMN price DROP NOT NULL
		`
		if err := dbs.DB.Exec(alterPriceSQL).Error; err != nil {
			logger.Warn(ctx, fmt.Sprintf("Price column might already be nullable: %v", err))
		}
	}

	logger.Info(ctx, "✓ Successfully migrated RoomTypeAdditional")
	return nil
}

func (dbs *DBPostgre) migrateBookingDetailAdditional(ctx context.Context) error {
	logger.Info(ctx, "Starting BookingDetailAdditional migration")

	tableName := "booking_detail_additionals"

	// Step 1: Check and add Category column
	var categoryExists bool
	checkCategorySQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'category'
		)
	`
	if err := dbs.DB.Raw(checkCategorySQL, tableName).Scan(&categoryExists).Error; err != nil {
		return fmt.Errorf("failed to check category column: %w", err)
	}

	if !categoryExists {
		logger.Info(ctx, "Adding category column to booking_detail_additionals")
		addCategorySQL := fmt.Sprintf(`
			ALTER TABLE booking_detail_additionals 
			ADD COLUMN category VARCHAR(10) DEFAULT '%s'
		`, constant.AdditionalServiceCategoryPrice)
		if err := dbs.DB.Exec(addCategorySQL).Error; err != nil {
			return fmt.Errorf("failed to add category column: %w", err)
		}
		// Set default value for existing records
		updateCategorySQL := fmt.Sprintf(`
			UPDATE booking_detail_additionals 
			SET category = '%s' 
			WHERE category IS NULL OR category = ''
		`, constant.AdditionalServiceCategoryPrice)
		if err := dbs.DB.Exec(updateCategorySQL).Error; err != nil {
			return fmt.Errorf("failed to set default category: %w", err)
		}
	}

	// Step 2: Check and add Pax column
	var paxExists bool
	checkPaxSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'pax'
		)
	`
	if err := dbs.DB.Raw(checkPaxSQL, tableName).Scan(&paxExists).Error; err != nil {
		return fmt.Errorf("failed to check pax column: %w", err)
	}

	if !paxExists {
		logger.Info(ctx, "Adding pax column to booking_detail_additionals")
		addPaxSQL := `
			ALTER TABLE booking_detail_additionals 
			ADD COLUMN pax INTEGER
		`
		if err := dbs.DB.Exec(addPaxSQL).Error; err != nil {
			return fmt.Errorf("failed to add pax column: %w", err)
		}
	}

	// Step 3: Check and add IsRequired column
	var isRequiredExists bool
	checkIsRequiredSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'is_required'
		)
	`
	if err := dbs.DB.Raw(checkIsRequiredSQL, tableName).Scan(&isRequiredExists).Error; err != nil {
		return fmt.Errorf("failed to check is_required column: %w", err)
	}

	if !isRequiredExists {
		logger.Info(ctx, "Adding is_required column to booking_detail_additionals")
		addIsRequiredSQL := `
			ALTER TABLE booking_detail_additionals 
			ADD COLUMN is_required BOOLEAN DEFAULT false
		`
		if err := dbs.DB.Exec(addIsRequiredSQL).Error; err != nil {
			return fmt.Errorf("failed to add is_required column: %w", err)
		}
	}

	// Step 4: Make Price column nullable
	var priceNullable bool
	checkPriceNullableSQL := `
		SELECT is_nullable FROM information_schema.columns 
		WHERE table_name = $1 AND column_name = 'price'
	`
	var isNullableStr string
	if err := dbs.DB.Raw(checkPriceNullableSQL, tableName).Scan(&isNullableStr).Error; err != nil {
		return fmt.Errorf("failed to check price column nullability: %w", err)
	}
	priceNullable = (isNullableStr == "YES")

	if !priceNullable {
		logger.Info(ctx, "Making price column nullable in booking_detail_additionals")
		alterPriceNullableSQL := `
			ALTER TABLE booking_detail_additionals 
			ALTER COLUMN price DROP NOT NULL
		`
		if err := dbs.DB.Exec(alterPriceNullableSQL).Error; err != nil {
			return fmt.Errorf("failed to make price column nullable: %w", err)
		}
	}

	// Step 5: Drop foreign key constraint to room_type_additionals to allow deleting room_type_additionals
	// without affecting historical booking_detail_additionals (they already store a full snapshot).
	logger.Info(ctx, "Dropping foreign key constraint from booking_detail_additionals to room_type_additionals if exists")
	dropFKSQL := `
		ALTER TABLE booking_detail_additionals 
		DROP CONSTRAINT IF EXISTS fk_booking_detail_additionals_room_type_additional
	`
	if err := dbs.DB.Exec(dropFKSQL).Error; err != nil {
		logger.Warn(ctx, fmt.Sprintf("Failed to drop foreign key constraint (might not exist): %v", err))
	}

	logger.Info(ctx, "✓ Successfully migrated BookingDetailAdditional")
	return nil
}

func (dbs *DBPostgre) migrateBookingGuest(ctx context.Context) error {
	logger.Info(ctx, "Starting BookingGuest migration")

	tableName := "booking_guests"

	// Step 1: Check and add Honorific column
	var honorificExists bool
	checkHonorificSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'honorific'
		)
	`
	if err := dbs.DB.Raw(checkHonorificSQL, tableName).Scan(&honorificExists).Error; err != nil {
		return fmt.Errorf("failed to check honorific column: %w", err)
	}

	if !honorificExists {
		logger.Info(ctx, "Adding honorific column to booking_guests")
		addHonorificSQL := `
			ALTER TABLE booking_guests 
			ADD COLUMN honorific VARCHAR(10)
		`
		if err := dbs.DB.Exec(addHonorificSQL).Error; err != nil {
			return fmt.Errorf("failed to add honorific column: %w", err)
		}
	}

	// Step 2: Check and add Category column
	var categoryExists bool
	checkCategorySQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'category'
		)
	`
	if err := dbs.DB.Raw(checkCategorySQL, tableName).Scan(&categoryExists).Error; err != nil {
		return fmt.Errorf("failed to check category column: %w", err)
	}

	if !categoryExists {
		logger.Info(ctx, "Adding category column to booking_guests")
		addCategorySQL := fmt.Sprintf(`
			ALTER TABLE booking_guests 
			ADD COLUMN category VARCHAR(20) DEFAULT '%s'
		`, constant.GuestCategoryAdult)
		if err := dbs.DB.Exec(addCategorySQL).Error; err != nil {
			return fmt.Errorf("failed to add category column: %w", err)
		}
		// Set default value for existing records
		updateCategorySQL := fmt.Sprintf(`
			UPDATE booking_guests 
			SET category = '%s' 
			WHERE category IS NULL OR category = ''
		`, constant.GuestCategoryAdult)
		if err := dbs.DB.Exec(updateCategorySQL).Error; err != nil {
			return fmt.Errorf("failed to set default category: %w", err)
		}
	}

	// Step 3: Check and add Age column
	var ageExists bool
	checkAgeSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'age'
		)
	`
	if err := dbs.DB.Raw(checkAgeSQL, tableName).Scan(&ageExists).Error; err != nil {
		return fmt.Errorf("failed to check age column: %w", err)
	}

	if !ageExists {
		logger.Info(ctx, "Adding age column to booking_guests")
		addAgeSQL := `
			ALTER TABLE booking_guests 
			ADD COLUMN age INTEGER
		`
		if err := dbs.DB.Exec(addAgeSQL).Error; err != nil {
			return fmt.Errorf("failed to add age column: %w", err)
		}
	}

	logger.Info(ctx, "✓ Successfully migrated BookingGuest")
	return nil
}

func (dbs *DBPostgre) migrateBookingDetailBedType(ctx context.Context) error {
	logger.Info(ctx, "Starting BookingDetail bed_type migration")

	tableName := "booking_details"

	// Check and add BedType column
	var bedTypeExists bool
	checkBedTypeSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'bed_type'
		)
	`
	if err := dbs.DB.Raw(checkBedTypeSQL, tableName).Scan(&bedTypeExists).Error; err != nil {
		return fmt.Errorf("failed to check bed_type column: %w", err)
	}

	if !bedTypeExists {
		logger.Info(ctx, "Adding bed_type column to booking_details")
		addBedTypeSQL := `
			ALTER TABLE booking_details 
			ADD COLUMN bed_type TEXT
		`
		if err := dbs.DB.Exec(addBedTypeSQL).Error; err != nil {
			return fmt.Errorf("failed to add bed_type column: %w", err)
		}
	}

	logger.Info(ctx, "✓ Successfully migrated BookingDetail bed_type")
	return nil
}

func (dbs *DBPostgre) migrateBookingDetailAdditionalNotes(ctx context.Context) error {
	logger.Info(ctx, "Starting BookingDetail additional_notes migration")

	tableName := "booking_details"

	// Check and add AdditionalNotes column
	var additionalNotesExists bool
	checkAdditionalNotesSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'additional_notes'
		)
	`
	if err := dbs.DB.Raw(checkAdditionalNotesSQL, tableName).Scan(&additionalNotesExists).Error; err != nil {
		return fmt.Errorf("failed to check additional_notes column: %w", err)
	}

	if !additionalNotesExists {
		logger.Info(ctx, "Adding additional_notes column to booking_details")
		addAdditionalNotesSQL := `
			ALTER TABLE booking_details 
			ADD COLUMN additional_notes TEXT
		`
		if err := dbs.DB.Exec(addAdditionalNotesSQL).Error; err != nil {
			return fmt.Errorf("failed to add additional_notes column: %w", err)
		}
	}

	logger.Info(ctx, "✓ Successfully migrated BookingDetail additional_notes")
	return nil
}

func (dbs *DBPostgre) migrateBookingDetailAdminNotes(ctx context.Context) error {
	logger.Info(ctx, "Starting BookingDetail admin_notes migration")

	tableName := "booking_details"

	// Check and add AdminNotes column
	var adminNotesExists bool
	checkAdminNotesSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'admin_notes'
		)
	`
	if err := dbs.DB.Raw(checkAdminNotesSQL, tableName).Scan(&adminNotesExists).Error; err != nil {
		return fmt.Errorf("failed to check admin_notes column: %w", err)
	}

	if !adminNotesExists {
		logger.Info(ctx, "Adding admin_notes column to booking_details")
		addAdminNotesSQL := `
			ALTER TABLE booking_details 
			ADD COLUMN admin_notes TEXT
		`
		if err := dbs.DB.Exec(addAdminNotesSQL).Error; err != nil {
			return fmt.Errorf("failed to add admin_notes column: %w", err)
		}
	}

	logger.Info(ctx, "✓ Successfully migrated BookingDetail admin_notes")
	return nil
}

// ✅ FUNGSI BARU: Migrasi Multi-Currency support
func (dbs *DBPostgre) migrateMultiCurrency(ctx context.Context) error {
	logger.Info(ctx, "Starting Multi-Currency migration")

	// Step 1: Add Prices JSONB column to room_prices
	if err := dbs.migrateRoomPricesJSONB(ctx); err != nil {
		return fmt.Errorf("failed to migrate room_prices prices: %w", err)
	}

	// Step 2: Add Prices JSONB column to room_type_additionals
	if err := dbs.migrateRoomTypeAdditionalsPricesJSONB(ctx); err != nil {
		return fmt.Errorf("failed to migrate room_type_additionals prices: %w", err)
	}

	// Step 3: Add Currency column to users
	if err := dbs.migrateUsersCurrency(ctx); err != nil {
		return fmt.Errorf("failed to migrate users currency: %w", err)
	}

	// Step 4: Add Currency column to booking_details
	if err := dbs.migrateBookingDetailsCurrency(ctx); err != nil {
		return fmt.Errorf("failed to migrate booking_details currency: %w", err)
	}

	// Step 5: Migrate promos.detail structure to support per-currency prices
	if err := dbs.migratePromoDetailStructure(ctx); err != nil {
		return fmt.Errorf("failed to migrate promo detail structure: %w", err)
	}

	// Step 6: Create GIN indexes on JSONB columns
	if err := dbs.createCurrencyIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create currency indexes: %w", err)
	}

	logger.Info(ctx, "✓ Successfully migrated Multi-Currency support")
	return nil
}

// Migrate room_prices: Add Prices JSONB and migrate existing Price to Prices
func (dbs *DBPostgre) migrateRoomPricesJSONB(ctx context.Context) error {
	logger.Info(ctx, "Migrating room_prices to support multi-currency")

	tableName := "room_prices"

	// Check if Prices column already exists
	var pricesExists bool
	checkPricesSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'prices'
		)
	`
	if err := dbs.DB.Raw(checkPricesSQL, tableName).Scan(&pricesExists).Error; err != nil {
		return fmt.Errorf("failed to check prices column: %w", err)
	}

	if !pricesExists {
		logger.Info(ctx, "Adding prices JSONB column to room_prices")
		addPricesSQL := `
			ALTER TABLE room_prices 
			ADD COLUMN prices JSONB
		`
		if err := dbs.DB.Exec(addPricesSQL).Error; err != nil {
			return fmt.Errorf("failed to add prices column: %w", err)
		}
	}

	// Migrate existing Price to Prices: {"IDR": price}
	logger.Info(ctx, "Migrating existing price data to prices JSONB")
	migratePriceSQL := `
		UPDATE room_prices 
		SET prices = jsonb_build_object('IDR', price)
		WHERE prices IS NULL OR prices = '{}'::jsonb
		AND price IS NOT NULL AND price > 0
	`
	if err := dbs.DB.Exec(migratePriceSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate price to prices: %w", err)
	}

	logger.Info(ctx, "✓ Successfully migrated room_prices prices")
	return nil
}

// Migrate room_type_additionals: Add Prices JSONB and migrate existing Price to Prices
func (dbs *DBPostgre) migrateRoomTypeAdditionalsPricesJSONB(ctx context.Context) error {
	logger.Info(ctx, "Migrating room_type_additionals to support multi-currency")

	tableName := "room_type_additionals"

	// Check if Prices column already exists
	var pricesExists bool
	checkPricesSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'prices'
		)
	`
	if err := dbs.DB.Raw(checkPricesSQL, tableName).Scan(&pricesExists).Error; err != nil {
		return fmt.Errorf("failed to check prices column: %w", err)
	}

	if !pricesExists {
		logger.Info(ctx, "Adding prices JSONB column to room_type_additionals")
		addPricesSQL := `
			ALTER TABLE room_type_additionals 
			ADD COLUMN prices JSONB
		`
		if err := dbs.DB.Exec(addPricesSQL).Error; err != nil {
			return fmt.Errorf("failed to add prices column: %w", err)
		}
	}

	// Migrate existing Price to Prices: {"IDR": price} (only for category="price" and price is not null)
	logger.Info(ctx, "Migrating existing price data to prices JSONB")
	migratePriceSQL := `
		UPDATE room_type_additionals 
		SET prices = jsonb_build_object('IDR', price)
		WHERE prices IS NULL OR prices = '{}'::jsonb
		AND category = 'price'
		AND price IS NOT NULL AND price > 0
	`
	if err := dbs.DB.Exec(migratePriceSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate price to prices: %w", err)
	}

	logger.Info(ctx, "✓ Successfully migrated room_type_additionals prices")
	return nil
}

// Migrate users: Add Currency column with default 'IDR'
func (dbs *DBPostgre) migrateUsersCurrency(ctx context.Context) error {
	logger.Info(ctx, "Migrating users to support currency preference")

	tableName := "users"

	// Check if Currency column already exists
	var currencyExists bool
	checkCurrencySQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'currency'
		)
	`
	if err := dbs.DB.Raw(checkCurrencySQL, tableName).Scan(&currencyExists).Error; err != nil {
		return fmt.Errorf("failed to check currency column: %w", err)
	}

	if !currencyExists {
		logger.Info(ctx, "Adding currency column to users")
		addCurrencySQL := `
			ALTER TABLE users 
			ADD COLUMN currency VARCHAR(3) DEFAULT 'IDR'
		`
		if err := dbs.DB.Exec(addCurrencySQL).Error; err != nil {
			return fmt.Errorf("failed to add currency column: %w", err)
		}
	}

	// Set default 'IDR' for existing records
	logger.Info(ctx, "Setting default currency 'IDR' for existing users")
	updateCurrencySQL := `
		UPDATE users 
		SET currency = 'IDR' 
		WHERE currency IS NULL OR currency = ''
	`
	if err := dbs.DB.Exec(updateCurrencySQL).Error; err != nil {
		return fmt.Errorf("failed to set default currency: %w", err)
	}

	logger.Info(ctx, "✓ Successfully migrated users currency")
	return nil
}

// Migrate booking_details: Add Currency column with default 'IDR'
func (dbs *DBPostgre) migrateBookingDetailsCurrency(ctx context.Context) error {
	logger.Info(ctx, "Migrating booking_details to support currency")

	tableName := "booking_details"

	// Check if Currency column already exists
	var currencyExists bool
	checkCurrencySQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'currency'
		)
	`
	if err := dbs.DB.Raw(checkCurrencySQL, tableName).Scan(&currencyExists).Error; err != nil {
		return fmt.Errorf("failed to check currency column: %w", err)
	}

	if !currencyExists {
		logger.Info(ctx, "Adding currency column to booking_details")
		addCurrencySQL := `
			ALTER TABLE booking_details 
			ADD COLUMN currency VARCHAR(3) DEFAULT 'IDR'
		`
		if err := dbs.DB.Exec(addCurrencySQL).Error; err != nil {
			return fmt.Errorf("failed to add currency column: %w", err)
		}
	}

	// Set default 'IDR' for existing records
	logger.Info(ctx, "Setting default currency 'IDR' for existing booking_details")
	updateCurrencySQL := `
		UPDATE booking_details 
		SET currency = 'IDR' 
		WHERE currency IS NULL OR currency = ''
	`
	if err := dbs.DB.Exec(updateCurrencySQL).Error; err != nil {
		return fmt.Errorf("failed to set default currency: %w", err)
	}

	logger.Info(ctx, "✓ Successfully migrated booking_details currency")
	return nil
}

// Create GIN indexes on JSONB columns for better query performance
func (dbs *DBPostgre) createCurrencyIndexes(ctx context.Context) error {
	logger.Info(ctx, "Creating GIN indexes on JSONB price columns")

	// Index on room_prices.prices
	indexName1 := "idx_room_prices_prices"
	checkIndex1SQL := `
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes 
			WHERE indexname = $1
		)
	`
	var index1Exists bool
	if err := dbs.DB.Raw(checkIndex1SQL, indexName1).Scan(&index1Exists).Error; err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}

	if !index1Exists {
		logger.Info(ctx, "Creating GIN index on room_prices.prices")
		createIndex1SQL := `
			CREATE INDEX idx_room_prices_prices ON room_prices USING GIN (prices)
		`
		if err := dbs.DB.Exec(createIndex1SQL).Error; err != nil {
			return fmt.Errorf("failed to create index on room_prices.prices: %w", err)
		}
	}

	// Index on room_type_additionals.prices
	indexName2 := "idx_room_type_additionals_prices"
	var index2Exists bool
	if err := dbs.DB.Raw(checkIndex1SQL, indexName2).Scan(&index2Exists).Error; err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}

	if !index2Exists {
		logger.Info(ctx, "Creating GIN index on room_type_additionals.prices")
		createIndex2SQL := `
			CREATE INDEX idx_room_type_additionals_prices ON room_type_additionals USING GIN (prices)
		`
		if err := dbs.DB.Exec(createIndex2SQL).Error; err != nil {
			return fmt.Errorf("failed to create index on room_type_additionals.prices: %w", err)
		}
	}

	logger.Info(ctx, "✓ Successfully created currency indexes")
	return nil
}

// Migrate promos.detail JSONB structure to support per-currency prices
// Old structure: {"fixed_price": 100, "discount_percentage": 10}
// New structure: {"prices": {"IDR": 1200000, "USD": 75}, "discount_percentage": 10}
func (dbs *DBPostgre) migratePromoDetailStructure(ctx context.Context) error {
	logger.Info(ctx, "Migrating promos.detail structure to support multi-currency")

	tableName := "promos"

	// Check if detail column exists (it should, but let's be safe)
	var detailExists bool
	checkDetailSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = $1 AND column_name = 'detail'
		)
	`
	if err := dbs.DB.Raw(checkDetailSQL, tableName).Scan(&detailExists).Error; err != nil {
		return fmt.Errorf("failed to check detail column: %w", err)
	}

	if !detailExists {
		logger.Warn(ctx, "Detail column does not exist in promos table, skipping migration")
		return nil
	}

	// Migrate existing fixed_price to prices structure
	// For promos with fixed_price, add prices: {"IDR": fixed_price}
	// Keep fixed_price field for backward compatibility during transition
	logger.Info(ctx, "Migrating existing fixed_price to prices structure in promo details")
	migratePromoDetailSQL := `
		UPDATE promos
		SET detail = jsonb_set(
			detail,
			'{prices}',
			jsonb_build_object('IDR', (detail->>'fixed_price')::float)
		)
		WHERE detail ? 'fixed_price'
		AND (detail->>'fixed_price') IS NOT NULL
		AND (detail->>'fixed_price')::float > 0
		AND (detail->'prices' IS NULL OR detail->'prices' = '{}'::jsonb)
	`
	if err := dbs.DB.Exec(migratePromoDetailSQL).Error; err != nil {
		return fmt.Errorf("failed to migrate promo detail structure: %w", err)
	}

	logger.Info(ctx, "✓ Successfully migrated promo detail structure")
	return nil
}

// ✅ FUNGSI BARU: Migrasi Currency Symbol - Update IDR symbol from "Rp" to "IDR"
func (dbs *DBPostgre) migrateCurrencySymbol(ctx context.Context) error {
	logger.Info(ctx, "Starting Currency Symbol migration: Updating IDR symbol from 'Rp' to 'IDR'")

	tableName := "currencies"

	// Check if currencies table exists
	var tableExists bool
	checkTableSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = $1
		)
	`
	if err := dbs.DB.Raw(checkTableSQL, tableName).Scan(&tableExists).Error; err != nil {
		return fmt.Errorf("failed to check currencies table existence: %w", err)
	}

	if !tableExists {
		logger.Info(ctx, "Currencies table does not exist, skipping migration")
		return nil
	}

	// Check if IDR currency exists with old symbol "Rp"
	var idrCount int64
	checkIDRSQL := `
		SELECT COUNT(*) FROM currencies 
		WHERE code = 'IDR' AND symbol = 'Rp'
	`
	if err := dbs.DB.Raw(checkIDRSQL).Scan(&idrCount).Error; err != nil {
		return fmt.Errorf("failed to check IDR currency: %w", err)
	}

	if idrCount == 0 {
		logger.Info(ctx, "No IDR currency found with symbol 'Rp', migration not needed")
		return nil
	}

	// Update IDR currency symbol from "Rp" to "IDR"
	logger.Info(ctx, fmt.Sprintf("Updating %d IDR currency record(s) from symbol 'Rp' to 'IDR'", idrCount))
	updateSymbolSQL := `
		UPDATE currencies 
		SET symbol = 'IDR'
		WHERE code = 'IDR' AND symbol = 'Rp'
	`
	if err := dbs.DB.Exec(updateSymbolSQL).Error; err != nil {
		return fmt.Errorf("failed to update IDR currency symbol: %w", err)
	}

	// Verify the update
	var updatedCount int64
	verifySQL := `
		SELECT COUNT(*) FROM currencies 
		WHERE code = 'IDR' AND symbol = 'IDR'
	`
	if err := dbs.DB.Raw(verifySQL).Scan(&updatedCount).Error; err != nil {
		return fmt.Errorf("failed to verify currency symbol update: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("✓ Successfully migrated Currency Symbol: %d IDR record(s) updated to symbol 'IDR'", updatedCount))
	return nil
}

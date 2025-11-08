package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
	"wtm-backend/config"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/infrastructure/database/seed"
	"wtm-backend/pkg/logger"
)

type DBPostgre struct {
	DB *gorm.DB
}

type dbContextKey struct{}

func NewDBPostgre(cfg *config.Config) (*DBPostgre, error) {
	ctx := context.Background()
	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable connect_timeout=5",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresName,
		cfg.PostgresPort,
	)

	logger.Info(ctx, "Connecting to PostgreSQL...", dataSourceName)
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", err.Error())
		return nil, err
	}

	// Optimize connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", err.Error())
		return nil, err
	}

	sqlDB.SetMaxOpenConns(100)                 // Maximum open connections (100 recommended for high traffic)
	sqlDB.SetMaxIdleConns(50)                  // Maximum idle connections (50 for faster reuse)
	sqlDB.SetConnMaxLifetime(10 * time.Minute) // Connection max lifetime to prevent memory leaks
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Close idle connections after 5 minutes

	logger.Info(ctx, "Database connected successfully with optimized pooling")

	// Auto-migrate tables
	err = db.AutoMigrate(
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
	)
	if err != nil {
		logger.Error(ctx, "Error during database migration", err.Error())
		return nil, err
	}

	logger.Info(ctx, "Database migration completed")

	seed.Seeding(db)

	return &DBPostgre{DB: db}, nil
}

func (dbs *DBPostgre) ErrRecordNotFound(ctx context.Context, err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warn(ctx,
			"Record not found")
		return true
	}

	return false
}

func (dbs *DBPostgre) ErrDuplicateKey(ctx context.Context, err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // unique_violation
			logger.Warn(ctx, "Duplicate key violation: %s", pgErr.ConstraintName)
			return true
		}
	}
	return false
}

func (dbs *DBPostgre) BeginTrx(ctx context.Context) (*gorm.DB, context.Context, error) {
	tx := dbs.DB.Begin()
	if tx.Error != nil {
		logger.Error(ctx, "Failed to begin transaction", tx.Error)
		dbs.resetDBSession()
		return tx, ctx, tx.Error
	}
	txCtx := context.WithValue(ctx, dbContextKey{}, tx)
	return tx, txCtx, nil
}

func (dbs *DBPostgre) CommitTrx(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		logger.Error(ctx, "Transaction is nil, cannot commit")
		return errors.New("transaction is nil")
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error(ctx, "Failed to commit transaction", err.Error())
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			logger.Error(ctx, "Failed to rollback transaction after commit failure", rollbackErr.Error())
			return fmt.Errorf("commit failed and rollback also failed: %s", rollbackErr.Error())
		}
		dbs.resetDBSession()
		return err
	}

	logger.Info(ctx, "Transaction committed successfully")
	dbs.resetDBSession()
	return nil
}

func (dbs *DBPostgre) RollbackTrx(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		logger.Error(ctx, "Transaction is nil, cannot rollback")
		return errors.New("transaction is nil")
	}

	if err := tx.Rollback().Error; err != nil {
		logger.Error(ctx, "Failed to rollback transaction", err.Error())
		dbs.resetDBSession()
		return err
	}

	logger.Info(ctx, "Transaction rolled back successfully")
	dbs.resetDBSession()
	return nil
}

func (dbs *DBPostgre) GetTx(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(dbContextKey{}).(*gorm.DB)
	if ok && db != nil {
		return db
	}
	if dbs.DB == nil {
		logger.Error(ctx, "Database connection is nil in context")
		return nil
	}
	return dbs.DB
}

func (dbs *DBPostgre) resetDBSession() {
	dbs.DB = dbs.DB.Session(&gorm.Session{NewDB: true})
}

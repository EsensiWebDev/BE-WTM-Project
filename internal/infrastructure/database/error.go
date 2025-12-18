package database

import (
	"context"
	"errors"
	"wtm-backend/pkg/logger"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (dbs *DBPostgre) ErrRecordNotFound(ctx context.Context, err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warn(ctx, "Record not found")
		return true
	}
	return false
}

// IsRecordNotFound checks if error is ErrRecordNotFound without logging.
// Use this when "not found" is expected behavior (e.g., checking if cart exists).
func (dbs *DBPostgre) IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (dbs *DBPostgre) ErrDuplicateKey(ctx context.Context, err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		logger.Warn(ctx, "Duplicate key violation", pgErr.ConstraintName)
		return true
	}
	return false
}

func (dbs *DBPostgre) ErrForeignKeyViolation(ctx context.Context, err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		logger.Warn(ctx, "Foreign key violation", pgErr.ConstraintName)
		return true
	}
	return false
}

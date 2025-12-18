package driver

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type DBPostgre interface {
	// Transaction Management
	BeginTrx(ctx context.Context) (*gorm.DB, context.Context, error)
	CommitTrx(ctx context.Context, tx *gorm.DB) error
	RollbackTrx(ctx context.Context, tx *gorm.DB) error
	GetTx(ctx context.Context) *gorm.DB
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error

	// Error Handlers
	ErrRecordNotFound(ctx context.Context, err error) bool
	IsRecordNotFound(err error) bool // Silent check for expected "not found" cases
	ErrDuplicateKey(ctx context.Context, err error) bool
	ErrForeignKeyViolation(ctx context.Context, err error) bool

	// Health & Monitoring
	HealthCheck(ctx context.Context) error
	GetStats(ctx context.Context) (*sql.DBStats, error)
	Close() error

	// Direct DB access (optional - untuk legacy code)
	GetDB() *gorm.DB
}

package driver

import (
	"context"
	"gorm.io/gorm"
)

type DBPostgre interface {
	BeginTrx(ctx context.Context) (*gorm.DB, context.Context, error)
	CommitTrx(ctx context.Context, tx *gorm.DB) error
	RollbackTrx(ctx context.Context, tx *gorm.DB) error
	GetTx(ctx context.Context) *gorm.DB
	ErrRecordNotFound(ctx context.Context, err error) bool
	ErrDuplicateKey(ctx context.Context, rr error) bool
}

package driver

import (
	"context"
	"wtm-backend/internal/infrastructure/database"
	"wtm-backend/pkg/logger"
)

type DatabaseTransaction struct {
	db *database.DBPostgre
}

func NewDatabaseTransaction(db *database.DBPostgre) *DatabaseTransaction {
	return &DatabaseTransaction{db: db}
}

func (dr *DatabaseTransaction) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, txCtx, err := dr.db.BeginTrx(ctx)
	if err != nil {
		logger.Error(ctx, "Error starting transaction", err.Error())
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = dr.db.RollbackTrx(txCtx, tx)
			panic(r)
		}
	}()

	if err = fn(txCtx); err != nil {
		logger.Error(txCtx, "Error executing transaction", err.Error())
		if rollbackErr := dr.db.RollbackTrx(txCtx, tx); rollbackErr != nil {
			logger.Error(txCtx, "Error rolling back transaction", rollbackErr.Error())
			return rollbackErr
		}
		return err
	}

	if err = dr.db.CommitTrx(txCtx, tx); err != nil {
		logger.Error(txCtx, "Error committing transaction", err.Error())
		_ = dr.db.RollbackTrx(txCtx, tx) // fallback rollback
		return err
	}

	return nil
}

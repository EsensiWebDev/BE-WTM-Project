package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"wtm-backend/pkg/logger"
)

// Key untuk menyimpan transaction di context
type dbContextKey struct{}

func (dbs *DBPostgre) BeginTrx(ctx context.Context) (*gorm.DB, context.Context, error) {
	if dbs.DB == nil {
		return nil, ctx, errors.New("database connection is nil")
	}

	tx := dbs.DB.Begin()
	if tx.Error != nil {
		logger.Error(ctx, "Begin transaction failed", tx.Error.Error())
		return nil, ctx, fmt.Errorf("begin transaction: %w", tx.Error)
	}

	txCtx := context.WithValue(ctx, dbContextKey{}, tx)
	return tx, txCtx, nil
}

func (dbs *DBPostgre) CommitTrx(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error(ctx, "Commit transaction failed", err.Error())
		dbs.rollbackWithLog(ctx, tx, "after commit failure")
		return fmt.Errorf("commit: %w", err)
	}

	logger.Info(ctx, "Transaction committed")
	dbs.resetDBSession()
	return nil
}

func (dbs *DBPostgre) RollbackTrx(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	if err := tx.Rollback().Error; err != nil {
		logger.Error(ctx, "Rollback transaction failed", err.Error())
		dbs.resetDBSession()
		return fmt.Errorf("rollback: %w", err)
	}

	logger.Info(ctx, "Transaction rolled back")
	dbs.resetDBSession()
	return nil
}

func (dbs *DBPostgre) GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(dbContextKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return dbs.DB
}

func (dbs *DBPostgre) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, txCtx, err := dbs.BeginTrx(ctx)
	if err != nil {
		return err
	}

	var operationErr error
	func() {
		defer func() {
			if p := recover(); p != nil {
				dbs.RollbackTrx(ctx, tx)
				panic(p)
			}
		}()
		operationErr = fn(txCtx)
	}()

	if operationErr != nil {
		if rbErr := dbs.RollbackTrx(ctx, tx); rbErr != nil {
			return fmt.Errorf("operation: %w, rollback: %v", operationErr, rbErr)
		}
		return operationErr
	}

	if err := dbs.CommitTrx(ctx, tx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (dbs *DBPostgre) HealthCheck(ctx context.Context) error {
	if dbs.DB == nil {
		return errors.New("database connection is nil")
	}

	sqlDB, err := dbs.DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	return sqlDB.Ping()
}

func (dbs *DBPostgre) GetStats(ctx context.Context) (*sql.DBStats, error) {
	sqlDB, err := dbs.DB.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return &stats, nil
}

func (dbs *DBPostgre) Close() error {
	if dbs.DB == nil {
		return nil
	}

	sqlDB, err := dbs.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (dbs *DBPostgre) GetDB() *gorm.DB {
	return dbs.DB
}

// Private methods
func (dbs *DBPostgre) resetDBSession() {
	if dbs.DB != nil && dbs.baseConfig != nil {
		dbs.DB = dbs.DB.Session(&gorm.Session{
			NewDB:                  true,
			PrepareStmt:            dbs.baseConfig.PrepareStmt,
			SkipDefaultTransaction: dbs.baseConfig.SkipDefaultTransaction,
		})
	}
}

func (dbs *DBPostgre) rollbackWithLog(ctx context.Context, tx *gorm.DB, reason string) {
	if err := tx.Rollback().Error; err != nil {
		logger.Error(ctx, "Rollback "+reason+" failed", err.Error())
	}
	dbs.resetDBSession()
}

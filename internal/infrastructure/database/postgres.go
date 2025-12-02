package database

import (
	"context"
	"fmt"
	"wtm-backend/config"
	"wtm-backend/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBPostgre(cfg *config.Config) (*DBPostgre, error) {
	return NewDBPostgreWithConfig(cfg, DefaultDBConfig)
}

func NewDBPostgreWithConfig(cfg *config.Config, dbConfig DBConfig) (*DBPostgre, error) {
	ctx := context.Background()

	// Security: jangan log credentials
	logger.Info(ctx, "Connecting to PostgreSQL",
		fmt.Sprintf("host=%s db=%s", cfg.PostgresHost, cfg.PostgresName))

	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable connect_timeout=%d",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresName,
		cfg.PostgresPort,
		dbConfig.ConnectTimeout,
	)

	baseConfig := &gorm.Config{
		PrepareStmt: true, // Performance improvement
	}

	db, err := gorm.Open(postgres.Open(dataSourceName), baseConfig)
	if err != nil {
		logger.Error(ctx, "Database connection failed", err.Error())
		return nil, fmt.Errorf("database connection: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	logger.Info(ctx, "Database connected",
		fmt.Sprintf("pool: %d/%d", dbConfig.MaxIdleConns, dbConfig.MaxOpenConns))

	dbs := &DBPostgre{
		DB:         db,
		baseConfig: baseConfig,
	}

	// Run migrations
	if err := dbs.runMigrations(ctx, cfg); err != nil {
		return nil, err
	}

	return dbs, nil
}

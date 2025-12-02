// database/types.go
package database

import (
	"gorm.io/gorm"
	"time"
)

type DBPostgre struct {
	DB         *gorm.DB
	baseConfig *gorm.Config
}

type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	ConnectTimeout  int
}

var DefaultDBConfig = DBConfig{
	MaxOpenConns:    100,
	MaxIdleConns:    50,
	ConnMaxLifetime: 10 * time.Minute,
	ConnMaxIdleTime: 5 * time.Minute,
	ConnectTimeout:  5,
}

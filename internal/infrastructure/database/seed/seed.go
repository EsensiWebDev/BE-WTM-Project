package seed

import (
	"context"
	"wtm-backend/pkg/logger"

	"gorm.io/gorm"
)

type Seed struct {
	db *gorm.DB
}

func NewSeed(db *gorm.DB) *Seed {
	return &Seed{
		db: db,
	}
}

func Seeding(db *gorm.DB) error {
	ctx := context.Background()
	logger.Info(ctx, "Start Seeding database...")

	seed := NewSeed(db)
	seed.SeedCurrency()
	seed.SeedUser()
	seed.SeedHotel()
	seed.SeedPromo()
	seed.SeedBooking()
	seed.SeedEmailTemplate()
	seed.SeedingEmailLog()

	logger.Info(ctx, "Finished Seeding database")

	return nil

}

package seed

import (
	"context"
	"gorm.io/gorm"
	"wtm-backend/pkg/logger"
)

type Seed struct {
	db *gorm.DB
}

func NewSeed(db *gorm.DB) *Seed {
	return &Seed{
		db: db,
	}
}

func Seeding(db *gorm.DB) {
	ctx := context.Background()
	logger.Info(ctx, "Start Seeding database...")

	seed := NewSeed(db)
	seed.SeedUser()
	seed.SeedHotel()
	seed.SeedPromo()
	seed.SeedBooking()
	seed.SeedEmailTemplate()
	seed.SeedingEmailLog()

	logger.Info(ctx, "Finished Seeding database")

}

package booking_repository

import (
	"context"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
	"wtm-backend/internal/domain"
	"wtm-backend/internal/infrastructure/database"
	"wtm-backend/pkg/logger"
)

type BookingRepository struct {
	db          *database.DBPostgre
	redisClient domain.RedisClient
}

func NewBookingRepository(db *database.DBPostgre, redisClient domain.RedisClient) *BookingRepository {
	return &BookingRepository{
		db:          db,
		redisClient: redisClient,
	}
}

func (br *BookingRepository) generateCode(ctx context.Context, keyRedis string, prefixCode string) (string, error) {
	date := time.Now().Format("060102") // yyMMdd
	redisKey := fmt.Sprintf("%s:%s", keyRedis, date)

	for i := 0; i < 10; i++ {
		suffix, _ := gonanoid.New(4)
		used, err := br.redisClient.IsSuffixUsed(ctx, redisKey, suffix)
		if err != nil {
			logger.Warn(ctx,
				"failed to check if redis is used", "error", err)
			return "", err
		}
		if !used {
			if err := br.redisClient.MarkSuffixUsed(ctx, redisKey, suffix); err != nil {
				logger.Warn(ctx, "failed to mark redis as used", "error", err)
				return "", err
			}
			return fmt.Sprintf("%s-%s-%s", prefixCode, date, suffix), nil
		}
	}

	return "", fmt.Errorf("failed to generate unique booking code after 10 attempts")
}

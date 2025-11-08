package domain

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	IsSuffixUsed(ctx context.Context, dateKey string, suffix string) (bool, error)
	MarkSuffixUsed(ctx context.Context, dateKey string, suffix string) error
	TTL(ctx context.Context, key string) (time.Duration, error)
}

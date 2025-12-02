package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"wtm-backend/config"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(config *config.Config) (*RedisClient, error) {
	ctx := context.Background()
	addressRedis := fmt.Sprintf("%s:%s", config.HostRedis, config.PortRedis)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addressRedis,
		Password: config.PasswordRedis,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Fatal(ctx, "Failed to initialize Redis client", err.Error())
		return nil, err
	}

	return &RedisClient{Client: rdb}, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.Client.Set(ctx, key, value, expiration).Err(); err != nil {
		logger.Error(ctx, "Error setting value in Redis", err.Error())
		return err
	}
	return nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.Warn(ctx, "Key not found in Redis", key)
			return "", nil
		}
		logger.Error(ctx, "Error getting value from Redis", err.Error())
		return "", err
	}

	return result, nil
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		logger.Error(ctx, "Error deleting key from Redis", err.Error())
		return err
	}

	return nil
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.Warn(ctx, "Key does not exist in Redis", key)
			return false, nil
		}
		logger.Error(ctx, "Error checking key existence in Redis", err.Error())
		return false, err
	}
	return count > 0, err
}

// IsSuffixUsed checks if a suffix already exists in the Redis Set for a given date
func (r *RedisClient) IsSuffixUsed(ctx context.Context, dateKey string, suffix string) (bool, error) {
	return r.Client.SIsMember(ctx, dateKey, suffix).Result()
}

// MarkSuffixUsed adds a suffix to the Redis Set and sets TTL if not already set
func (r *RedisClient) MarkSuffixUsed(ctx context.Context, dateKey string, suffix string) error {
	pipe := r.Client.TxPipeline()

	pipe.SAdd(ctx, dateKey, suffix)
	pipe.Expire(ctx, dateKey, getEndOfDayTTL())

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error(ctx, "Error marking suffix as used in Redis", err.Error())
	}
	return err
}

func getEndOfDayTTL() time.Duration {
	now := time.Now().In(constant.AsiaJakarta)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return time.Until(endOfDay)
}

func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

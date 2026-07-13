package services

import (
	"context"
	"errors"
	"time"

	"go-service/infrastructure/exceptions"
	"go-service/infrastructure/redis/interfaces"

	"github.com/redis/go-redis/v9"
)

type RedisCacheService struct {
	client *redis.Client
}

func NewRedisCacheService(client *redis.Client) interfaces.RedisCacheInterface {
	return &RedisCacheService{client: client}
}

func (r *RedisCacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		panic(exceptions.ServerErrorException(err))
	}
}

func (r *RedisCacheService) Get(ctx context.Context, key string) string {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return ""
	}
	if err != nil {
		panic(exceptions.ServerErrorException(err))
	}
	return val
}

func (r *RedisCacheService) Delete(ctx context.Context, key string) {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		panic(exceptions.ServerErrorException(err))
	}
}

func (r *RedisCacheService) Exists(ctx context.Context, key string) bool {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		panic(exceptions.ServerErrorException(err))
	}
	return result > 0
}

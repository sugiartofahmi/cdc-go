package interfaces

import (
	"context"
	"time"
)

type RedisCacheInterface interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration)
	Get(ctx context.Context, key string) string
	Delete(ctx context.Context, key string)
	Exists(ctx context.Context, key string) bool
}

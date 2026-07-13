package singleton

import "github.com/redis/go-redis/v9"

func RedisSingleton() *redis.Client {
	return redisClient
}

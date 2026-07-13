package singleton

import (
	"sync"

	"go-service/infrastructure/integrations"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	once        sync.Once
	httpClient  *integrations.HttpClient
	db          *gorm.DB
	redisClient *redis.Client
)

func Init(httpClientInstance *integrations.HttpClient, dbInstance *gorm.DB, redisClientInstance *redis.Client) {
	once.Do(func() {
		httpClient = httpClientInstance
		db = dbInstance
		redisClient = redisClientInstance
	})
}

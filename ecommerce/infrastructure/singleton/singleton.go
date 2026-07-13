package singleton

import (
	"sync"

	"go-service/infrastructure/integrations"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	once             sync.Once
	httpClient       *integrations.HttpClient
	db               *gorm.DB
	redisClient      *redis.Client
	opensearchClient *opensearchapi.Client
)

func Init(httpClientInstance *integrations.HttpClient, dbInstance *gorm.DB, redisClientInstance *redis.Client, opensearchClientInstance *opensearchapi.Client) {
	once.Do(func() {
		httpClient = httpClientInstance
		db = dbInstance
		redisClient = redisClientInstance
		opensearchClient = opensearchClientInstance
	})
}

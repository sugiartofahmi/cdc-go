package singleton

import (
	"sync"

	"go-service/infrastructure/integrations"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/redis/go-redis/v9"
)

var (
	once             sync.Once
	httpClient       *integrations.HttpClient
	redisClient      *redis.Client
	opensearchClient *opensearchapi.Client
)

func Init(httpClientInstance *integrations.HttpClient, redisClientInstance *redis.Client, opensearchClientInstance *opensearchapi.Client) {
	once.Do(func() {
		httpClient = httpClientInstance
		redisClient = redisClientInstance
		opensearchClient = opensearchClientInstance
	})
}

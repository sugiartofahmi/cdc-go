package singleton

import (
	"sync"

	"go-service/infrastructure/integrations"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
)

var (
	once             sync.Once
	httpClient       *integrations.HttpClient
	redisClient      *redis.Client
	opensearchClient *opensearchapi.Client
	kafka *kafkago.Reader
)

func Init(httpClientInstance *integrations.HttpClient, redisClientInstance *redis.Client, opensearchClientInstance *opensearchapi.Client, kafkaInstance *kafkago.Reader) {
	once.Do(func() {
		httpClient = httpClientInstance
		redisClient = redisClientInstance
		opensearchClient = opensearchClientInstance
		kafka = kafkaInstance
	})
}

package interfaces

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerInterface interface {
	Consume(ctx context.Context, handler func(msg kafka.Message) error)
	Close() error
}

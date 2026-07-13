package services

import (
	"context"
	"log"

	"go-service/infrastructure/exceptions"
	"go-service/infrastructure/kafka/interfaces"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumerService struct {
	reader *kafka.Reader
}

func NewKafkaConsumerService(reader *kafka.Reader) interfaces.KafkaConsumerInterface {
	return &KafkaConsumerService{reader: reader}
}

func (s *KafkaConsumerService) Consume(ctx context.Context, handler func(msg kafka.Message) error) {
	for {
		select {
		case <-ctx.Done():
			log.Println("kafka consumer stopped")
			return
		default:
			msg, err := s.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Println("error fetching kafka message:", err)
				continue
			}

			if err := handler(msg); err != nil {
				log.Println("error handling kafka message:", err)
			}

			if err := s.reader.CommitMessages(ctx, msg); err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Println("error committing kafka message:", err)
				panic(exceptions.ServerErrorException(err))
			}
		}
	}
}

func (s *KafkaConsumerService) Close() error {
	return s.reader.Close()
}

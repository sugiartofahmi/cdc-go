package factories

import (
	"log"
	"strings"

	"go-service/infrastructure/config"

	"github.com/segmentio/kafka-go"
)

func NewKafka() (*kafka.Reader, error) {
	brokers := strings.Split(config.KafkaBrokers, ",")

	conn := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     config.KafkaGroupID,
		GroupTopics: []string{"ecommerce.public.products", "ecommerce.public.categories"},
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	log.Println("kafka connection created successfully")
	return conn, nil
}

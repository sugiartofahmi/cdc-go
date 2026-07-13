package singleton

import kafkago "github.com/segmentio/kafka-go"

func KafkaSingleton() *kafkago.Reader {
	return kafka
}

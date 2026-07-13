package config

var (
	KafkaBrokers = Get("KAFKA_BROKERS", "localhost:9092")
	KafkaGroupID = Get("KAFKA_GROUP_ID", "ecommerce-search-group")
)

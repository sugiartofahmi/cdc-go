package config

var (
	RedisHost     = Get("REDIS_HOST", "localhost")
	RedisPort     = Get("REDIS_PORT", "6379")
	RedisPassword = Get("REDIS_PASS", "")
)

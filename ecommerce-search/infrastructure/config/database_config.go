package config

import (
	"os"
	"strconv"
)

var (
	DBHost                      = GetRequired("DB_HOST")
	DBPort                      = Get("DB_PORT", "5432")
	DBName                      = GetRequired("DB_NAME")
	DBUser                      = GetRequired("DB_USER")
	DBPassword                  = GetRequired("DB_PASSWORD")
	DBTimezone                  = Get("DB_TIMEZONE", "Asia/Jakarta")
	DBMaxIdleConns              = stringToInt(Get("DB_MAX_IDLE_CONNS", "5"))
	DBMaxOpenConns              = stringToInt(Get("DB_MAX_OPEN_CONNS", "10"))
	DBMaxIdleConnsInMinutes     = stringToInt(Get("DB_MAX_IDLE_CONNS_IN_MINUTES", "10"))
	DBMaxLifetimeConnsInMinutes = stringToInt(Get("DB_MAX_LIFETIME_CONNS_IN_MINUTES", "60"))
)

func stringToInt(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intValue
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

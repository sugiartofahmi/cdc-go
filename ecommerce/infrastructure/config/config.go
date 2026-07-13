package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

func LoadConfig() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Printf("warning: failed to load .env file: %v", err)
		}
	})
}

func Get(key, defaultValue string) string {
	LoadConfig()
	value := os.Getenv(key)
	isValueExist := value != ""
	if isValueExist {
		return value
	}
	return defaultValue
}

func GetRequired(key string) string {
	LoadConfig()
	value := os.Getenv(key)
	isValueExist := value != ""
	if !isValueExist {
		return key + "-not-set"
	}
	return value
}

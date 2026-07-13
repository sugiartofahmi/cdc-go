package config

var AppName = Get("APP_NAME", "go-service")
var AppEnv = Get("APP_ENV", "development")
var AppPort = Get("APP_PORT", "8080")
var AppGinMode = Get("APP_GIN_MODE", "debug")

package config

var (
	JWTSecret    = GetRequired("JWT_SECRET")
	JWTExpiredIn = Get("JWT_EXPIRES_IN", "24h")
)

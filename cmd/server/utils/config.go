package utils

import (
	"os"
)

type Config struct {
	AppEnv               string
	AppVersion           string
	AuthEnabled          bool
	BaseUrl              string
	DebugRedisEnabled    bool
	DebugRequestsEnabled bool
	DebugSqlEnabled      bool
	DevModeEnabled       bool
	HTTPPort             string
	PostgresDsn          string
	RedisDsn             string
}

func LoadConfig() Config {
	return Config{
		AppEnv:               GetEnv("APP_ENV", "dev"),
		AppVersion:           GetEnv("APP_VERSION", "0.0.0"),
		AuthEnabled:          GetEnv("APP_HTTP_AUTH_ENABLED", "0") == "1",
		BaseUrl:              GetEnv("APP_BASE_URL", "https://localhost:4444"),
		DebugRedisEnabled:    GetEnv("APP_DEBUG_REDIS_ENABLED", "0") == "1",
		DebugRequestsEnabled: GetEnv("APP_DEBUG_REQUESTS_ENABLED", "0") == "1",
		DebugSqlEnabled:      GetEnv("APP_DEBUG_SQL_ENABLED", "1") == "1",
		DevModeEnabled:       GetEnv("APP_ENV", "dev") == "dev",
		PostgresDsn:          GetEnv("APP_DATABASE_DSN", "postgres://postgres:example@127.0.0.1:5432/interview?sslmode=disable"),
		RedisDsn:             GetEnv("APP_REDIS_DSN", "localhost:6379"),
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

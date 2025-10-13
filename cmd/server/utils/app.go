package utils

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config   *Config
	Logger   *slog.Logger
	Postgres *pgxpool.Pool
	Redis    *redis.Client
}

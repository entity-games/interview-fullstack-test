package utils

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXQueryTracer struct {
	Logger *slog.Logger
}

var whitespaceRegex = regexp.MustCompile(`\s+`)

func (t *PGXQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	formattedArgs := make([]string, len(data.Args))
	for i, arg := range data.Args {
		s := fmt.Sprint(arg)
		if len(s) > 20 {
			s = s[:20] + "..."
		}
		formattedArgs[i] = s
	}

	t.Logger.InfoContext(ctx, "[SQL]",
		slog.String("sql", strings.TrimSpace(whitespaceRegex.ReplaceAllString(data.SQL, " "))),
		slog.String("args", strings.Join(formattedArgs, ", ")),
	)
	return ctx
}

func (t *PGXQueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
}

func InitPostgresPool(config Config) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(config.PostgresDsn)
	if err != nil {
		slog.Error("Failed to parse config: %v", "err", err)
		return nil
	}

	if config.DebugSqlEnabled {
		cfg.ConnConfig.Tracer = &PGXQueryTracer{
			Logger: slog.Default(),
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		slog.Error("Unable to connect to database: %v", "err", err)
		return nil
	}

	return pool
}

func PingPostgres(ctx context.Context, db *pgxpool.Pool) bool {
	conn, err := db.Acquire(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Connection error: %v", "err", err)
		return false
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT 1")
	if err != nil {
		slog.ErrorContext(ctx, "Query error: %v", "err", err)
		return false
	}
	defer rows.Close()

	if rows.Next() {
		var result int
		if err := rows.Scan(&result); err != nil {
			slog.ErrorContext(ctx, "Scan error: %v", "err", err)
			return false
		}
		return result == 1
	}

	return false
}

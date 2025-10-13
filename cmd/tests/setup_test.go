package tests

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"entity/interview/cmd/server/routing"
	"entity/interview/cmd/server/utils"
)

var suite struct {
	app     utils.App
	handler http.Handler
}

func bindToHandler(t *testing.T, handler http.Handler) *httpexpect.Expect {
	t.Helper()

	client := &http.Client{
		Transport: httpexpect.NewBinder(handler),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // stop after the first response
		},
		Jar: httpexpect.NewCookieJar(),
	}

	return httpexpect.WithConfig(httpexpect.Config{
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
		Printers: []httpexpect.Printer{}, // httpexpect.NewDebugPrinter(t, true)
	})
}

func flushPostgres(db *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.Exec(ctx, "TRUNCATE users")
	db.Exec(ctx, "TRUNCATE auth")
	db.Exec(ctx, "TRUNCATE game_data")
	db.Exec(ctx, "TRUNCATE orders")
}

func flushRedis(db *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.FlushDB(ctx)
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func SetupTest(t *testing.T) *httpexpect.Expect {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	flushPostgres(suite.app.Postgres)
	flushRedis(suite.app.Redis)

	return bindToHandler(t, suite.handler)
}

func TestMain(m *testing.M) {
	err := utils.InitialiseTemplates()
	if err != nil {
		slog.Error("Could not init templates", "err", err)
		panic(err)
	}

	config := utils.LoadConfig()
	config.DebugSqlEnabled = false // enable to see SQL queries for debugging failed tests

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	redisPool := utils.InitRedisPool(config)
	postgresPool := utils.InitPostgresPool(config)
	defer postgresPool.Close()

	app := utils.App{Config: &config, Logger: logger, Postgres: postgresPool, Redis: redisPool}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := routing.SetupGinEngine(r, app)

	suite.app = app
	suite.handler = handler

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	os.Exit(m.Run())
}

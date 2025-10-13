package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/routing"
	"entity/interview/cmd/server/utils"
)

func main() {
	err := utils.InitialiseTemplates()
	if err != nil {
		slog.Error("Could not init templates", "err", err)
		return
	}

	config := utils.LoadConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	postgresPool := utils.InitPostgresPool(config)
	defer postgresPool.Close()

	redisPool := utils.InitRedisPool(config)

	app := utils.App{Config: &config, Logger: logger, Postgres: postgresPool, Redis: redisPool}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	handler := routing.SetupGinEngine(r, app)

	slog.Info("SYSTEM Server running on http://localhost:4444")
	err = http.ListenAndServe(":4444", handler)
	slog.Error("err", err)
}

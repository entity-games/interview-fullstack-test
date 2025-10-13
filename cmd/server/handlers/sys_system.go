package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/utils"
)

func ServeHealthCheck(c *gin.Context, a utils.App) {
	c.JSON(http.StatusOK, gin.H{"version": a.Config.AppVersion})
}

func ServeConnectionCheck(c *gin.Context, a utils.App) {
	postgresConnected := utils.PingPostgres(c.Request.Context(), a.Postgres)
	redisConnected := utils.PingRedis(c.Request.Context(), a.Redis)

	c.JSON(http.StatusOK, gin.H{
		"version":  a.Config.AppVersion,
		"postgres": postgresConnected,
		"redis":    redisConnected,
	})
}

func ServeNotFound(c *gin.Context, a utils.App) {
	c.AbortWithStatus(http.StatusNotFound)
}

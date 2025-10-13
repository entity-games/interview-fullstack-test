package routing

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/handlers"
	"entity/interview/cmd/server/middleware"
	"entity/interview/cmd/server/utils"
)

func SetupGinEngine(r *gin.Engine, app utils.App) http.Handler {
	// Add Gin-style middleware, the order is important
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLoggingMiddleware("/ping"))
	r.Use(middleware.WithTimeoutMiddleware())

	web := r.Group("/")

	r.Use(cors.New(cors.Config{
		AllowOrigins:              []string{"*"},
		AllowMethods:              []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:              []string{"Content-OrderType", "Authorization"},
		AllowCredentials:          true,
		OptionsResponseStatusCode: http.StatusNoContent,
	}))

	platform := r.Group("/v1/platform")
	platform.Use(middleware.SessionAuthMiddleware(app.Redis))

	// System
	r.GET("/ping", func(c *gin.Context) { handlers.ServeHealthCheck(c, app) })
	r.GET("/status", func(c *gin.Context) { handlers.ServeConnectionCheck(c, app) })

	// Platform Public API, requiring Session ID, with /v1/platform prefix
	platform.POST("/makePurchase", func(c *gin.Context) { handlers.MakeVirtualPurchase(c, app) })
	platform.GET("/getCurrentUser", func(c *gin.Context) { handlers.GetCurrentUser(c, app) })

	// Static
	web.GET("/static/*path", func(c *gin.Context) { handlers.ServeStatic(c, app) })
	web.GET("/favicon.ico", func(c *gin.Context) { handlers.ServeFavicon(c, app) })

	// Web app
	web.POST("/login", func(c *gin.Context) { handlers.Login(c, app) })
	web.GET("/logout", func(c *gin.Context) { handlers.Logout(c, app) })
	web.GET("/game/:linkType/:gameId/:titleStub", func(c *gin.Context) { handlers.ServeGame(c, app) })
	web.GET("/partial/header", func(c *gin.Context) { handlers.ServeHeader(c, app) })

	// Default
	web.GET("/", func(c *gin.Context) { handlers.ServeHomepage(c, app) })
	web.GET("/index.html", func(c *gin.Context) { handlers.ServeHomepage(c, app) })

	r.NoRoute(func(c *gin.Context) { handlers.ServeNotFound(c, app) })

	return r.Handler()
}

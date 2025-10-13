package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func RequestLoggingMiddleware(skipPrefixes ...string) gin.HandlerFunc {
	inner := gin.Logger() // the standard Gin logger

	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(path, prefix) {
				c.Next()
				return
			}
		}
		inner(c)
	}
}

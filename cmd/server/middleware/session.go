package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"entity/interview/cmd/server/data"
)

const SessionContextKey = "session"

func SessionAuthMiddleware(db *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		sessionId := c.GetHeader("X-Session-Id")
		if sessionId == "" {
			slog.ErrorContext(ctx, "Missing session identifier when required")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "Invalid Session ID"})
			return
		}

		// Validate session
		session := data.SessionData{}
		err := data.SessionGet(ctx, db, &session, sessionId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "Invalid Session ID"})
			return
		}

		if session.UserID == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "Invalid Session ID"})
			return
		}

		c.Set(SessionContextKey, &session)
		c.Next()
	}
}

func GetAuthSession(c *gin.Context) *data.SessionData {
	res, _ := c.Get(SessionContextKey)
	return res.(*data.SessionData)
}

package data

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID string `redis:"user_id"`
}

func SessionGet(ctx context.Context, rdb *redis.Client, sessionData *SessionData, sessionId string) error {
	if sessionId == "" {
		return errors.New("sessionId is empty")
	}
	sessionKey := "sessions::" + sessionId

	err := rdb.HGetAll(ctx, sessionKey).Scan(sessionData)
	return err
}

func SessionSet(ctx context.Context, rdb *redis.Client, sessionId string, sessionData *SessionData) bool {
	if sessionId == "" {
		return false
	}
	sessionKey := "sessions::" + sessionId

	p := rdb.Pipeline()
	p.HMSet(ctx, sessionKey, sessionData)
	p.Expire(ctx, sessionKey, time.Hour*24*7)
	_, err := p.Exec(ctx)

	return err == nil
}

func SessionSetNewId(c *gin.Context) string {
	sessionId := uuid.New().String()

	c.SetCookie("session", sessionId, int(time.Hour.Seconds()*24), "/", "", true, false)
	return sessionId
}

func SessionReset(c *gin.Context, rdb *redis.Client) {
	sessionId, _ := c.Cookie("session")
	if sessionId == "" {
		return
	}
	rdb.Unlink(c.Request.Context(), "sessions::"+sessionId)
	SessionSetNewId(c)
}

func SessionReadIdFromCookie(c *gin.Context) string {
	sessionId, _ := c.Cookie("session")
	if sessionId == "" {
		sessionId = SessionSetNewId(c)
	}
	return sessionId
}

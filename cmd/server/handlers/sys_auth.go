package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/data"
	"entity/interview/cmd/server/utils"
)

func Logout(c *gin.Context, a utils.App) {
	data.SessionReset(c, a.Redis)
	c.Redirect(http.StatusSeeOther, "/")
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context, a utils.App) {
	ctx := c.Request.Context()

	req := LoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user := &data.LookupUser{Username: req.Username}
	err := data.UsersGetByUsername(ctx, a.Postgres, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid username"})
		return
	}

	// fake testing of passwords for the interview
	if req.Password != "fixturepass" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid password"})
		return
	}

	sessionId := data.SessionReadIdFromCookie(c)
	sessionData := &data.SessionData{
		UserID: user.UserID,
	}
	data.SessionSet(ctx, a.Redis, sessionId, sessionData)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

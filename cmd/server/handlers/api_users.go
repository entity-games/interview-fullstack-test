package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/data"
	"entity/interview/cmd/server/middleware"
	"entity/interview/cmd/server/utils"
)

/*** Requests & Responses ***/

type UserSearchRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
}

func GetCurrentUser(c *gin.Context, a utils.App) {
	ctx := c.Request.Context()

	session := middleware.GetAuthSession(c)
	user := &data.User{UserID: session.UserID}
	err := data.UsersGet(ctx, a.Postgres, user)
	if err != nil {
		slog.ErrorContext(ctx, "Error fetching user", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

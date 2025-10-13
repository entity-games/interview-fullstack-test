package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"entity/interview/cmd/server/data"
	"entity/interview/cmd/server/utils"
)

func ServeHomepage(c *gin.Context, a utils.App) {
	ctx := c.Request.Context()
	path := c.Request.URL.Path

	if path != "/" && path != "/index.html" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	sessionId := data.SessionReadIdFromCookie(c)
	session := &data.SessionData{}
	err := data.SessionGet(ctx, a.Redis, session, sessionId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve session: %s", err.Error()))
		return
	}

	var unlocked []int
	user := &data.User{}

	if session.UserID != "" {
		user.UserID = session.UserID
		err = data.UsersGet(ctx, a.Postgres, user)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve user: %s", err.Error()))
			return
		}

		unlocked = data.GameDataGetUnlocked(ctx, a.Postgres, session.UserID)
	}

	games := data.GamesList(ctx, a.Postgres, unlocked)
	details := map[string]any{
		"config":  a.Config,
		"session": session,
		"user":    user,
		"games":   games,
	}

	tmpl := utils.FindTemplate(c.Writer, "games.tmpl")
	if tmpl == nil {
		return
	}

	err = tmpl.Execute(c.Writer, details)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.ErrorContext(ctx, fmt.Sprintf("Failed executing template: %s", err.Error()))
	}
}

func ServeHeader(c *gin.Context, a utils.App) {
	ctx := c.Request.Context()

	sessionId := data.SessionReadIdFromCookie(c)
	session := &data.SessionData{}
	err := data.SessionGet(ctx, a.Redis, session, sessionId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve session: %s", err.Error()))
		return
	}

	var unlocked []int
	user := &data.User{}

	if session.UserID != "" {
		user.UserID = session.UserID
		err = data.UsersGet(ctx, a.Postgres, user)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Unable to retrieve user: %s", err.Error()))
			return
		}

		unlocked = data.GameDataGetUnlocked(ctx, a.Postgres, session.UserID)
	}

	games := data.GamesList(ctx, a.Postgres, unlocked)
	details := map[string]any{
		"config":  a.Config,
		"games":   games,
		"session": session,
		"user":    user,
	}

	tmpl := utils.FindTemplate(c.Writer, "header.tmpl")
	if tmpl == nil {
		return
	}

	err = tmpl.Execute(c.Writer, details)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.ErrorContext(ctx, fmt.Sprintf("Failed executing template: %s", err.Error()))
	}
}

func ServeGame(c *gin.Context, a utils.App) {
	ctx := c.Request.Context()

	var err error

	sessionId := data.SessionReadIdFromCookie(c)
	sessionData := &data.SessionData{}
	err = data.SessionGet(ctx, a.Redis, sessionData, sessionId)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid sessionData", "err", err)
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	gameId, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		slog.ErrorContext(ctx, "Invalid game", "err", err)
		c.String(http.StatusNotFound, "Invalid game")
		return
	}

	game := &data.Game{GameID: gameId}
	err = data.GamesGet(ctx, a.Postgres, game)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid game", "err", err)
		c.String(http.StatusNotFound, "Invalid game")
		return
	}

	tmpl := utils.FindTemplate(c.Writer, "game-canvas.tmpl")
	if tmpl == nil {
		return
	}

	details := map[string]any{
		"config":  a.Config,
		"session": sessionData,
	}

	err = tmpl.Execute(c.Writer, details)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		slog.ErrorContext(ctx, fmt.Sprintf("Failed executing template: %s", err.Error()))
	}
}

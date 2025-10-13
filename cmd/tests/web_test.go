package tests

import (
	"net/http"
	"testing"

	"github.com/google/uuid"

	"entity/interview/cmd/server/data"
)

func Test_Home_Page(t *testing.T) {
	e := SetupTest(t)

	body := e.GET("/").
		WithHeader("Accept", "text/html").
		Expect().
		Status(http.StatusOK).
		Body()

	body.Contains("Test Game")
}

func Test_Game_Page(t *testing.T) {
	e := SetupTest(t)

	e.GET("/game/s/111111/test-game").
		WithHeader("Accept", "text/html").
		Expect().
		Status(http.StatusOK)
}

func Test_Partial_Header(t *testing.T) {
	e := SetupTest(t)

	e.GET("/partial/header").
		WithHeader("Accept", "text/html").
		Expect().
		Status(http.StatusOK).
		Body().
		Contains("Login").
		Contains("Games").
		NotContains("Friends")

	// login
	user := &data.User{Username: "name_" + randomString(6)}
	_ = data.NewUserRegistration(t.Context(), suite.app.Postgres, user)

	sessionData := &data.SessionData{UserID: user.UserID}
	sessionId := uuid.New().String()
	data.SessionSet(t.Context(), suite.app.Redis, sessionId, sessionData)

	e.GET("/partial/header").
		WithHeader("Accept", "text/html").
		WithCookie("session", sessionId).
		Expect().
		Status(http.StatusOK).
		Body().
		NotContains("Login").
		Contains("Games").
		Contains("Friends")
}

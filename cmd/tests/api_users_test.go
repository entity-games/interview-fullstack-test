package tests

import (
	"net/http"
	"testing"

	"github.com/google/uuid"

	"entity/interview/cmd/server/data"
)

func Test_Users(t *testing.T) {
	e := SetupTest(t)

	user := &data.User{Username: "name_" + randomString(6)}
	_ = data.NewUserRegistration(t.Context(), suite.app.Postgres, user)

	sessionData := &data.SessionData{UserID: user.UserID}
	sessionId := uuid.New().String()
	data.SessionSet(t.Context(), suite.app.Redis, sessionId, sessionData)

	// fail getting the current user without a session
	e.GET("/v1/platform/getCurrentUser").
		WithHeader("Content-Type", "application/json").
		Expect().Status(http.StatusForbidden).
		JSON().Object().ContainsKey("err")

	e.GET("/v1/platform/getCurrentUser").
		WithHeader("Content-Type", "application/json").
		WithHeader("X-Session-Id", sessionId).
		Expect().Status(http.StatusOK).
		JSON().Object().
		HasValue("user_id", user.UserID).
		HasValue("username", user.Username).
		HasValue("coins", 50).
		ContainsKey("created_at").ContainsKey("updated_at")
}

func Test_Login(t *testing.T) {
	e := SetupTest(t)

	user := &data.User{Username: "name_" + randomString(6)}
	_ = data.NewUserRegistration(t.Context(), suite.app.Postgres, user)

	e.POST("/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(map[string]string{
			"username": user.Username,
			"password": "fixpass",
		}).
		Expect().Status(http.StatusOK).
		JSON().Object().Value("user").Object().HasValue("username", user.Username)
}

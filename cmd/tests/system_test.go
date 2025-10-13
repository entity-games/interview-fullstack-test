package tests

import (
	"net/http"
	"testing"
)

func Test_Ping(t *testing.T) {
	e := SetupTest(t)

	e.GET("/ping").
		WithHeader("Content-Type", "application/json").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("version")
}

func Test_Status(t *testing.T) {
	e := SetupTest(t)

	e.GET("/status").
		WithHeader("Content-Type", "application/json").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsKey("version")
}

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/labstack/echo/v4"
)

type MockHandler struct{}

func NewMockHandler() handler {
	return &MockHandler{}
}

func (mh *MockHandler) HandleGetMain(c echo.Context) error {
	Data := models.NewIndexData()

	time.Sleep(2 * time.Second)
	fmt.Println(c.Request().Cookies())
	fmt.Println(c.Cookie("sessionId"))

	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "testing!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Response(), &cookie)

	return c.Render(http.StatusOK, "index", Data)
}

func (mh *MockHandler) HandlePostProfile(c echo.Context) error {
	time.Sleep(2 * time.Second)
	return c.Render(http.StatusOK, "rooms", models.NewUser("anonymous", "/static/images/avatar_placeholder.png"))
}

func (mh *MockHandler) HandlePostRoom(c echo.Context) error {
	return nil
}

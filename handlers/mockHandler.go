package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MockHandler struct{}

func NewMockHandler() handler {
	return &MockHandler{}
}

func (mh *MockHandler) HandleGetMainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Main Page!")
}

func (mh *MockHandler) HandleGetRoomsPage(c echo.Context) error {
	return c.String(http.StatusOK, "Rooms Page!")
}

func (mh *MockHandler) HandlePostProfile(c echo.Context) error {
	return c.String(http.StatusOK, "Profle Posted!")
}

func (mh *MockHandler) HandlePutProfile(c echo.Context) error {
	return c.String(http.StatusOK, "Profile Updated!!")
}

func (mh *MockHandler) HandlePostRoom(c echo.Context) error {
	return c.String(http.StatusOK, "Room Created!")
}

func (mh *MockHandler) HandlePostChatMessage(c echo.Context) error {
	return c.String(http.StatusOK, "Message Posted!")
}

func (mh *MockHandler) HandleGetWebSocketConn(c echo.Context) error {
	return c.String(http.StatusOK, "Websocket Established!")
}

func (mh *MockHandler) HandleGetRoomChat(c echo.Context) error {
	return c.String(http.StatusOK, "Room Chat!")
}

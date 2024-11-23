package handlers

import (
	"net/http"
	"strconv"

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
	userIdString := c.Param("id")
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Getting rooms for user: "+strconv.Itoa(userId))
}

func (mh *MockHandler) HandlePostProfile(c echo.Context) error {
	return nil
}

func (mh *MockHandler) HandlePutProfile(c echo.Context) error {
	return nil
}

func (mh *MockHandler) HandlePostRoom(c echo.Context) error {
	return nil
}

func (mh *MockHandler) HandlePostChatMessage(c echo.Context) error {
	return nil
}

func (mh *MockHandler) HandleGetWebSocketConn(c echo.Context) error {
	return nil
}

func (mh *MockHandler) HandleGetRoomChat(c echo.Context) error {
	return nil
}

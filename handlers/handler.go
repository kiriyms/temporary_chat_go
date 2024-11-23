package handlers

import "github.com/labstack/echo/v4"

type handler interface {
	HandleGetMainPage(c echo.Context) error
	HandleGetRoomsPage(c echo.Context) error
	HandlePostProfile(c echo.Context) error
	HandlePutProfile(c echo.Context) error
	HandlePostRoom(c echo.Context) error
	HandleGetRoomChat(c echo.Context) error
	HandlePostChatMessage(c echo.Context) error
	HandleGetWebSocketConn(c echo.Context) error
}

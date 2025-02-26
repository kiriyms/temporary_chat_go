package handlers

import "github.com/labstack/echo/v4"

type Handler interface {
	HandleGetMain(c echo.Context) error
	HandlePostProfile(c echo.Context) error
	HandlePostRoom(c echo.Context) error
	HandleGetWebSocket(c echo.Context) error
	HandleGetWebSocketChat(c echo.Context) error
	HandleGetRoom(c echo.Context) error
	HandlePostJoinRoom(c echo.Context) error
	HandleDeleteLeaveRoom(c echo.Context) error
	HandleGetUserEditModal(c echo.Context) error
	HandlePostUserEdit(c echo.Context) error
}

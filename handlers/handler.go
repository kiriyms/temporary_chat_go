package handlers

import "github.com/labstack/echo/v4"

type Handler interface {
	HandleGetMain(c echo.Context) error
	HandlePostProfile(c echo.Context) error
	HandlePostRoom(c echo.Context) error
	HandleGetWebSocket(c echo.Context) error
	HandleGetRoom(c echo.Context) error
}

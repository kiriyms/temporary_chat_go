package handlers

import "github.com/labstack/echo/v4"

type handler interface {
	HandleGetMain(c echo.Context) error
	HandleGetRooms(c echo.Context) error
	HandlePostProfile(c echo.Context) error
}

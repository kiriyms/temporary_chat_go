package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProductionHandler struct{}

func NewProductionHandler() handler {
	return &ProductionHandler{}
}

func (ph *ProductionHandler) HandleGetMainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (ph *ProductionHandler) HandleGetRoomsPage(c echo.Context) error {
	return c.String(http.StatusOK, "Rooms Page!")
}

func (ph *ProductionHandler) HandlePostProfile(c echo.Context) error {
	return c.String(http.StatusOK, "Profle Posted!")
}

func (ph *ProductionHandler) HandlePutProfile(c echo.Context) error {
	return c.String(http.StatusOK, "Profile Updated!!")
}

func (ph *ProductionHandler) HandlePostRoom(c echo.Context) error {
	return c.String(http.StatusOK, "Room Created!")
}

func (ph *ProductionHandler) HandlePostChatMessage(c echo.Context) error {
	return c.String(http.StatusOK, "Message Posted!")
}

func (ph *ProductionHandler) HandleGetWebSocketConn(c echo.Context) error {
	return c.String(http.StatusOK, "Websocket Established!")
}

func (ph *ProductionHandler) HandleGetRoomChat(c echo.Context) error {
	return c.String(http.StatusOK, "Room Chat!")
}

func (ph *ProductionHandler) HandlePostToken(c echo.Context) error {
	return c.String(http.StatusOK, "Post Token!")
}

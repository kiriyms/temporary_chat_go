package main

import (
	"github.com/Kirill-Sirotkin/temporary_chat_go/handlers"
	"github.com/Kirill-Sirotkin/temporary_chat_go/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Renderer = utils.NewTemplates()
	e.Use(middleware.Logger())
	e.Static("/styles", "styles")

	h := handlers.NewProductionHandler()

	e.GET("/", h.HandleGetMainPage)
	e.POST("/", h.HandlePostProfile)
	e.PUT("/", h.HandlePutProfile)

	e.GET("/rooms", h.HandleGetRoomsPage)
	e.POST("/rooms", h.HandlePostRoom)

	e.GET("/rooms/:roomId", h.HandleGetRoomChat)
	e.POST("/rooms/:roomId", h.HandlePostChatMessage)
	e.GET("/rooms/:roomId/ws", h.HandleGetWebSocketConn)

	e.Logger.Fatal(e.Start(":1323"))
}

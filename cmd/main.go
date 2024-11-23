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
	h := handlers.NewMockHandler()

	e.GET("/", h.HandleGetMainPage)
	e.POST("/", h.HandlePostProfile)
	e.PUT("/:id", h.HandlePutProfile)

	e.GET("/:id/rooms", h.HandleGetRoomsPage)
	e.POST("/:id/rooms", h.HandlePostRoom)

	e.GET("/:id/rooms/:roomId", h.HandleGetRoomChat)
	e.POST("/:id/rooms/:roomId", h.HandlePostChatMessage)
	e.GET("/:id/rooms/:roomId/ws", h.HandleGetWebSocketConn)

	e.Logger.Fatal(e.Start(":1323"))
}

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
	e.Static("/static", "static")

	h := handlers.NewProductionHandler()

	e.GET("/", h.HandleGetMain)
	e.POST("/", h.HandlePostProfile)
	e.GET("/rooms", h.HandleGetRooms)

	e.Logger.Fatal(e.Start(":1323"))
}

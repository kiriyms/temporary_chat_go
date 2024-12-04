package main

import (
	"log"

	"github.com/Kirill-Sirotkin/temporary_chat_go/handlers"
	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/Kirill-Sirotkin/temporary_chat_go/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Renderer = utils.NewTemplates()
	e.Use(middleware.Logger())
	e.Static("/static", "static")

	userList := models.NewUserList()
	h := handlers.NewProductionHandler(userList)

	e.GET("/", h.HandleGetMain)
	e.POST("/", h.HandlePostProfile)
	e.GET("/rooms", h.HandleGetRooms)

	e.Logger.Fatal(e.Start(":1323"))
}

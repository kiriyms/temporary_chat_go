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
	e.Static("/uploads", "uploads")

	userList := models.NewUserList()
	roomList := models.NewRoomList()
	h := handlers.NewProductionHandler(userList, roomList)
	// h := handlers.NewMockHandler()

	e.GET("/", h.HandleGetMain)
	e.POST("/", h.HandlePostProfile)
	e.POST("/room", h.HandlePostRoom)

	e.Logger.Fatal(e.Start(":1323"))
}

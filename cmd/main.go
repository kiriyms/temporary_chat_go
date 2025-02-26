package main

import (
	"log"
	"strings"

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

	loggerCfg := middleware.LoggerConfig{
		Format:           `${time_custom} [ECHO]: ${method} "${host}" "${uri}"` + "\n" + strings.Repeat(" ", 28) + `status:${status}, error:"${error}", time:"${latency_human}"` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}

	e.Use(middleware.LoggerWithConfig(loggerCfg))
	e.Static("/static", "static")
	e.Static("/uploads", "uploads")

	userList := models.NewUserList()
	roomList := models.NewRoomList()
	var h handlers.Handler = handlers.NewProductionHandler(userList, roomList)
	// h := handlers.NewMockHandler()

	e.GET("/", h.HandleGetMain)
	e.POST("/", h.HandlePostProfile)
	e.POST("/room", h.HandlePostRoom)
	e.GET("/room/:roomId", h.HandleGetRoom)
	e.GET("/ws/:roomId", h.HandleGetWebSocket)
	e.GET("/ws/chat/:roomId", h.HandleGetWebSocketChat)
	e.POST("/room/join", h.HandlePostJoinRoom)
	e.DELETE("/:roomId", h.HandleDeleteLeaveRoom)
	e.GET("/edit", h.HandleGetUserEditModal)
	e.POST("/edit", h.HandlePostUserEdit)

	e.Logger.Fatal(e.Start(":1323"))
}

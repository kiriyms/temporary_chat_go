package main

import (
	"log"
	"os"
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
		log.Println("Error loading .env file, will look in environment variables")
	}

	if _, ok := os.LookupEnv("JWT_SECRET"); !ok {
		log.Fatal("JWT_SECRET is not set in environment variables")
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

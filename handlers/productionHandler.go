package handlers

import (
	"fmt"
	"net/http"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/Kirill-Sirotkin/temporary_chat_go/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ProductionHandler struct {
	UserList *models.UserList
	RoomList *models.RoomList
}

func NewProductionHandler(ul *models.UserList, rl *models.RoomList) *ProductionHandler {
	return &ProductionHandler{
		UserList: ul,
		RoomList: rl,
	}
}

func (ph *ProductionHandler) HandleGetMain(c echo.Context) error {
	data := models.NewIndexData()

	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		return c.Render(http.StatusOK, "index", data)
	}

	user := ph.UserList.GetUserById(userUUID)
	if user == nil {
		// error: user not found although token is valid
		fmt.Println("user not found although token is valid")
		return c.Render(http.StatusOK, "index", data)
	}

	data.Rooms = true
	data.UserData = user
	data.UserRooms = ph.RoomList.GetUserRooms(user.Id)
	data.CurrentRooms = len(ph.RoomList.GetUserRooms(user.Id))

	return c.Render(http.StatusOK, "index", data)
}

func (ph *ProductionHandler) HandlePostProfile(c echo.Context) error {
	userName := c.FormValue("name-input")
	fmt.Println(userName)

	file, err := c.FormFile("avatar-input")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error with file upload: "+err.Error())
	}
	fileName, err := utils.UploadFile(file)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error with file upload: "+err.Error())
	}
	fmt.Println(fileName)

	user := models.NewUser(userName, fileName)
	ph.UserList.AddUser(user)

	newUser := ph.UserList.GetUserById(user.Id)
	if newUser == nil {
		// handle error: no user in memory list error 500
		fmt.Println("User was not found!!!")
		return c.String(http.StatusInternalServerError, "User was not found!")
	}

	token, err := utils.CreateJWT(newUser.Id)
	if err != nil {
		fmt.Println("Could not create JWT!!!")
		return c.String(http.StatusInternalServerError, "Could not create JWT!")
	}

	jwtCookie := new(http.Cookie)
	jwtCookie.Name = "jwt"
	jwtCookie.Value = token
	c.SetCookie(jwtCookie)

	data := models.NewIndexData()
	data.Rooms = true
	data.UserData = newUser
	data.UserRooms = ph.RoomList.GetUserRooms(newUser.Id)

	return c.Render(http.StatusOK, "rooms", data)
}

func (ph *ProductionHandler) HandlePostRoom(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		// POST in htmx template should do an oob swap into id="main"
		// change http.Status to appropriate error
		return c.Render(http.StatusUnauthorized, "login", nil)
	}
	fmt.Println(userUUID)

	// besides the room list in-memory (possibly to be removed) needs to start a websocket(???)
	// or goroutine with a 3-5 minute timer. by the end the room is deleted
	// frontend html needs to have some js to show a timer and remove room element on-client

	hub := models.NewHub()
	go hub.Start()
	// room := models.NewRoom(userUUID, nil)
	room := models.NewRoom(userUUID, hub)
	hub.Id = room.Id
	ph.RoomList.AddRoom(room)
	utils.StartRoomTimer(ph.RoomList, room.Id)

	// refresh token
	token, err := utils.CreateJWT(userUUID)
	if err != nil {
		fmt.Println("Could not create JWT!!!")
		return c.String(http.StatusInternalServerError, "Could not create JWT!")
	}

	jwtCookie := new(http.Cookie)
	jwtCookie.Name = "jwt"
	jwtCookie.Value = token
	c.SetCookie(jwtCookie)

	dataCounter := struct {
		CurrentRooms int
		MaxRooms     int
	}{
		CurrentRooms: len(ph.RoomList.GetUserRooms(userUUID)),
		MaxRooms:     5,
	}

	data := struct {
		Id uuid.UUID
	}{
		room.Id,
	}

	c.Render(http.StatusOK, "room-list-counter-oob", dataCounter)
	return c.Render(http.StatusOK, "room-card", data)
}

func (ph *ProductionHandler) HandleGetWebSocket(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		// POST in htmx template should do an oob swap into id="main"
		// change http.Status to appropriate error
		return c.Render(http.StatusUnauthorized, "login", nil)
	}
	fmt.Println(userUUID)

	userRooms := ph.RoomList.GetUserRooms(userUUID)
	roomIdParam := c.Param("roomId")
	roomUUID, err := uuid.Parse(roomIdParam)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var hub *models.Hub = nil
	for _, room := range userRooms {
		if room.Id == roomUUID {
			hub = room.Hub
		}
	}
	if hub == nil {
		return c.String(http.StatusInternalServerError, "no room with that ID")
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	client := models.NewClient(userUUID, hub, conn)
	client.Hub.Register <- client

	go client.WriteToWebSocket()
	go client.ReadFromWebSocket()

	return nil
}

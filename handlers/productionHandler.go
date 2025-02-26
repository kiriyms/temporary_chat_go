package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
		log.Printf("[INFO]: token not found/invalid on initial page load. Serving index page with login. ErrMsg: %v", err)
		return c.Render(http.StatusOK, "index", data)
	}

	user := ph.UserList.GetUserById(userUUID)
	if user == nil {
		log.Printf("[ERROR]: user not found although token is valid")
		// notificationData := models.Notification{
		// 	IsError: true,
		// 	Content: "user not found although token is valid",
		// }
		// c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusOK, "index", data)
	}

	userRooms := ph.RoomList.GetUserRooms(user.Id)
	userRoomsWithTimer := make([]*models.RoomWithTimer, 0)

	for _, room := range userRooms {
		userRoomsWithTimer = append(userRoomsWithTimer, &models.RoomWithTimer{
			Room:         *room,
			TimerSeconds: int(time.Until(room.ExpireTime).Seconds()),
		})

		log.Printf("TIMER SECOND IN HANDLEGETMAIN: %v", int(time.Until(room.ExpireTime).Seconds()))
	}

	data.Rooms = true
	data.UserData = user
	data.UserRooms = userRoomsWithTimer
	data.CurrentRooms = len(ph.RoomList.GetUserRooms(user.Id))

	log.Printf("[INFO]: token found and valid. Serving index page with room data")
	return c.Render(http.StatusOK, "index", data)
}

func (ph *ProductionHandler) HandlePostProfile(c echo.Context) error {
	userName := c.FormValue("name-input")
	if userName == "" {
		userName = "Anonymous"
	}

	file, err := c.FormFile("avatar-input")
	fileName := "static/images/avatar_placeholder.png"
	if err == nil {
		fileName, err = utils.UploadFile(file)
		if err != nil {
			log.Printf("[ERROR]: error with file upload. ErrMsg: %v", err)
			return c.String(http.StatusInternalServerError, "Error with file upload: "+err.Error())
		}
	}

	user := models.NewUser(userName, fileName)
	ph.UserList.AddUser(user)

	newUser := ph.UserList.GetUserById(user.Id)
	if newUser == nil {
		log.Printf("[ERROR]: created user was not found in UserList")
		return c.String(http.StatusInternalServerError, "User was not found!")
	}
	resetChan := utils.NewUserTimer(ph.UserList, newUser.Id)
	newUser.ResetTimer = resetChan

	token, err := utils.CreateJWT(newUser.Id)
	if err != nil {
		log.Printf("[ERROR]: could not create JWT. ErrMsg: %v", err)
		return c.String(http.StatusInternalServerError, "Could not create JWT!")
	}

	jwtCookie := new(http.Cookie)
	jwtCookie.Name = "jwt"
	jwtCookie.Value = token
	jwtCookie.Path = "/"
	c.SetCookie(jwtCookie)

	userRooms := ph.RoomList.GetUserRooms(user.Id)
	userRoomsWithTimer := make([]*models.RoomWithTimer, 0)

	for _, room := range userRooms {
		userRoomsWithTimer = append(userRoomsWithTimer, &models.RoomWithTimer{
			Room:         *room,
			TimerSeconds: int(time.Until(room.ExpireTime).Seconds()),
		})
	}

	data := models.NewIndexData()
	data.Rooms = true
	data.UserData = newUser
	data.UserRooms = userRoomsWithTimer

	log.Printf("[INFO]: new profile posted. Uploaded name: %v, avatar filename: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"user info: id: %v, name: %v, avatarPath: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"active user rooms count: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"generated token: %v",
		userName, fileName, newUser.Id, newUser.Name, newUser.AvatarPath, len(data.UserRooms), token)
	return c.Render(http.StatusOK, "rooms", data)
}

func (ph *ProductionHandler) HandlePostRoom(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: POST room unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	roomName := c.FormValue("room-name")
	if roomName == "" {
		roomName = "New Room"
	}

	if len(ph.RoomList.GetUserRooms(userUUID)) >= 5 {
		dataCounter := struct {
			CurrentRooms int
			MaxRooms     int
		}{
			CurrentRooms: len(ph.RoomList.GetUserRooms(userUUID)),
			MaxRooms:     5,
		}
		return c.Render(http.StatusOK, "room-list-counter-error-oob", dataCounter)
	}

	roomCode := utils.GenerateRoomCode(6)
	for {
		if ph.RoomList.IsRoomCodeUsed(roomCode) {
			roomCode = utils.GenerateRoomCode(6)
			continue
		}
		break
	}

	hub := models.NewHub(ph.UserList)
	defer func() {
		go hub.Start()
	}()
	room := models.NewRoom(userUUID, hub, roomName, roomCode)
	utils.NewRoomTimer(ph.RoomList, room.Id)
	hub.Id = room.Id
	ph.RoomList.AddRoom(room)

	// refresh token
	token, err := utils.CreateJWT(userUUID)
	if err != nil {
		log.Printf("[ERROR]: could not create JWT. ErrMsg: %v", err)
		return c.String(http.StatusInternalServerError, "Could not create JWT!")
	}

	jwtCookie := new(http.Cookie)
	jwtCookie.Name = "jwt"
	jwtCookie.Value = token
	jwtCookie.Path = "/"
	c.SetCookie(jwtCookie)

	user := ph.UserList.GetUserById(userUUID)
	if user == nil {
		log.Printf("[ERROR]: user with id %v not found in userList", userUUID)
		return c.String(http.StatusInternalServerError, "User not found")
	}
	user.ResetTimer <- true

	dataCounter := struct {
		CurrentRooms int
		MaxRooms     int
	}{
		CurrentRooms: len(ph.RoomList.GetUserRooms(userUUID)),
		MaxRooms:     5,
	}

	data := struct {
		TimerSeconds int
		Room         struct {
			Id   uuid.UUID
			Name string
		}
	}{
		TimerSeconds: 60,
		Room: struct {
			Id   uuid.UUID
			Name string
		}{
			Id:   room.Id,
			Name: room.Name,
		},
	}

	log.Printf("[INFO]: new room posted. Uploaded room name: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"user info: id: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"room info: id: %v, name: %v, code: %v, userCount: %v, expireTime: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"hub info: id: %v, msgCount: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"new user JWT: %v",
		roomName, userUUID, room.Id, room.Name, room.Code, len(room.Users), room.ExpireTime, hub.Id, len(hub.Messages), token)
	c.Render(http.StatusOK, "room-list-new-room-input-oob", nil)
	c.Render(http.StatusOK, "room-list-counter-oob", dataCounter)
	return c.Render(http.StatusOK, "room-card", data)
}

// FIX "return err" to return c.Render with proper http status codes
func (ph *ProductionHandler) HandleGetWebSocket(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: GET websocket unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	// userRooms := ph.RoomList.GetUserRooms(userUUID)
	roomIdParam := c.Param("roomId")
	roomUUID, err := uuid.Parse(roomIdParam)
	if err != nil {
		log.Printf("[ERROR]: provided room id parameter '%v' could not be parsed into UUID. ErrMsg: %v", roomIdParam, err)
		return c.String(http.StatusInternalServerError, "Provided room id parameter could not be parsed into UUID")
	}

	var hub *models.Hub = nil
	for _, room := range ph.RoomList.GetRooms() {
		if room.Id == roomUUID {
			hub = room.Hub
		}
	}
	if hub == nil {
		log.Printf("[ERROR]: hub with id %v not found", roomUUID)
		return c.String(http.StatusInternalServerError, "No room with that ID")
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("[ERROR]: could not upgrade request to websocket. ErrMsg: %v", err)
		return c.String(http.StatusInternalServerError, "Could not upgrade request to websocket")
	}

	client := models.NewClient(userUUID, hub, conn)
	client.Hub.Register <- client

	defer func() {
		go client.WriteToWebSocket()
		go client.ReadFromWebSocket(false)
	}()

	log.Printf("[INFO]: upgraded request to websocket for room card with id %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"hub info: id: %v, msgCount: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"user info: id: %v",
		roomUUID, hub.Id, len(hub.Messages), userUUID)
	return nil
}

// FIX "return err" to return c.Render with proper http status codes
func (ph *ProductionHandler) HandleGetRoom(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: GET room unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	roomIdParam := c.Param("roomId")
	roomUUID, err := uuid.Parse(roomIdParam)
	if err != nil {
		log.Printf("[ERROR]: provided room id parameter '%v' could not be parsed into UUID. ErrMsg: %v", roomIdParam, err)
		return c.String(http.StatusInternalServerError, "Provided room id parameter could not be parsed into UUID")
	}

	if ph.RoomList.GetRoomById(roomUUID) == nil {
		log.Printf("[ERROR]: room with id %v not found", roomUUID)
		return c.String(http.StatusInternalServerError, "No room with that ID")
	}

	userRooms := ph.RoomList.GetUserRooms(userUUID)

	for _, room := range userRooms {
		roomData := struct {
			Id   uuid.UUID
			Name string
		}{
			Id:   room.Id,
			Name: room.Name,
		}
		if room.Id == roomUUID {
			c.Render(http.StatusOK, "room-card-active-oob", roomData)
			continue
		}
		c.Render(http.StatusOK, "room-card-inactive-oob", roomData)
	}

	roomUsers := ph.RoomList.GetRoomUsers(roomUUID)
	roomUsersInfo := make([]struct {
		Id         uuid.UUID
		Name       string
		AvatarPath string
	}, 0)
	for _, user := range roomUsers {
		userInfo := ph.UserList.GetUserById(user)
		roomUsersInfo = append(roomUsersInfo, struct {
			Id         uuid.UUID
			Name       string
			AvatarPath string
		}{
			Id:         userInfo.Id,
			Name:       userInfo.Name,
			AvatarPath: userInfo.AvatarPath,
		})
	}
	log.Printf("[INFO]: room users: %v", roomUsersInfo)
	room := ph.RoomList.GetRoomById(roomUUID)
	data := struct {
		Id        uuid.UUID
		Name      string
		Code      string
		RoomUsers []struct {
			Id         uuid.UUID
			Name       string
			AvatarPath string
		}
	}{
		Id:        room.Id,
		Name:      room.Name,
		Code:      room.Code,
		RoomUsers: roomUsersInfo,
	}

	log.Printf("[INFO]: sent room data"+
		"\n"+
		strings.Repeat(" ", 28)+
		"room info: id: %v, name: %v, code: %v, userCount: %v, expireTime: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"user info: id: %v",
		room.Id, room.Name, room.Code, len(room.Users), room.ExpireTime, userUUID)
	return c.Render(http.StatusOK, "chat-window-active", data)
}

func (ph *ProductionHandler) HandleGetWebSocketChat(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: GET websocket chat unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	// userRooms := ph.RoomList.GetUserRooms(userUUID)
	roomIdParam := c.Param("roomId")
	roomUUID, err := uuid.Parse(roomIdParam)
	if err != nil {
		log.Printf("[ERROR]: provided room id parameter '%v' could not be parsed into UUID. ErrMsg: %v", roomIdParam, err)
		return c.String(http.StatusInternalServerError, "Provided room id parameter could not be parsed into UUID")
	}

	var hub *models.Hub = nil
	for _, room := range ph.RoomList.GetRooms() {
		if room.Id == roomUUID {
			hub = room.Hub
		}
	}
	if hub == nil {
		log.Printf("[ERROR]: hub with id %v not found", roomUUID)
		return c.String(http.StatusInternalServerError, "No room with that ID")
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("[ERROR]: could not upgrade request to websocket. ErrMsg: %v", err)
		return c.String(http.StatusInternalServerError, "Could not upgrade request to websocket")
	}

	client := models.NewClient(userUUID, hub, conn)
	client.Hub.RegisterChat <- client

	defer func() {
		go client.WriteToWebSocket()
		go client.ReadFromWebSocket(true)
	}()

	log.Printf("[INFO]: upgraded request to websocket for chat with id %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"hub info: id: %v, msgCount: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"user info: id: %v",
		roomUUID, hub.Id, len(hub.Messages), userUUID)
	return nil
}

func (ph *ProductionHandler) HandlePostJoinRoom(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: POST join unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	roomCode := c.FormValue("room-code")
	room := ph.RoomList.GetRoomByCode(roomCode)

	if room == nil {
		log.Printf("[ERROR]: room with code %v not found", roomCode)
		return c.String(http.StatusInternalServerError, "No room with that code")
	}

	userRooms := ph.RoomList.GetUserRooms(userUUID)
	for _, room := range userRooms {
		if room.Code == roomCode {
			log.Printf("[ERROR]: user %v already in room %v", userUUID, roomCode)
			notificationData := models.Notification{
				IsError: true,
				Content: "you are already in this room",
			}
			return c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		}
	}

	// refresh token
	token, err := utils.CreateJWT(userUUID)
	if err != nil {
		log.Printf("[ERROR]: could not create JWT. ErrMsg: %v", err)
		return c.String(http.StatusInternalServerError, "Could not create JWT!")
	}

	if len(ph.RoomList.GetUserRooms(userUUID)) >= 5 {
		dataCounter := struct {
			CurrentRooms int
			MaxRooms     int
		}{
			CurrentRooms: len(ph.RoomList.GetUserRooms(userUUID)),
			MaxRooms:     5,
		}
		return c.Render(http.StatusOK, "room-list-counter-error-oob", dataCounter)
	}

	ph.RoomList.AddUserToRoom(room.Id, userUUID)

	jwtCookie := new(http.Cookie)
	jwtCookie.Name = "jwt"
	jwtCookie.Value = token
	jwtCookie.Path = "/"
	c.SetCookie(jwtCookie)

	user := ph.UserList.GetUserById(userUUID)
	if user == nil {
		log.Printf("[ERROR]: user with id %v not found in userList", userUUID)
		return c.String(http.StatusInternalServerError, "User not found")
	}

	user.ResetTimer <- true
	dataCounter := struct {
		CurrentRooms int
		MaxRooms     int
	}{
		CurrentRooms: len(ph.RoomList.GetUserRooms(userUUID)),
		MaxRooms:     5,
	}

	data := struct {
		TimerSeconds int
		Room         struct {
			Id   uuid.UUID
			Name string
		}
	}{
		TimerSeconds: int(time.Until(room.ExpireTime).Seconds()),
		Room: struct {
			Id   uuid.UUID
			Name string
		}{
			Id:   room.Id,
			Name: room.Name,
		},
	}

	log.Printf("[INFO]: joined user to room. Uploaded room code: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"room info: id: %v, name: %v, code: %v, userCount: %v, expireTime: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"seconds remaining in room: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"user info: id: %v"+
		"\n"+
		strings.Repeat(" ", 28)+
		"generated token: %v",
		roomCode, room.Id, room.Name, room.Code, len(room.Users), room.ExpireTime, int(time.Until(room.ExpireTime).Seconds()), userUUID, token)
	c.Render(http.StatusOK, "room-list-join-room-input-oob", nil)
	c.Render(http.StatusOK, "room-list-counter-oob", dataCounter)
	return c.Render(http.StatusOK, "room-card", data)
}

// FIX "return err" to return c.Render with proper http status codes
func (ph *ProductionHandler) HandleDeleteLeaveRoom(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: GET room unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	roomIdParam := c.Param("roomId")
	roomUUID, err := uuid.Parse(roomIdParam)
	if err != nil {
		log.Printf("[ERROR]: provided room id parameter '%v' could not be parsed into UUID. ErrMsg: %v", roomIdParam, err)
		return c.String(http.StatusInternalServerError, "Provided room id parameter could not be parsed into UUID")
	}

	err = ph.RoomList.RemoveUserFromRoom(roomUUID, userUUID)
	if err != nil {
		log.Printf("[ERROR]: %v", err)
		return c.String(http.StatusInternalServerError, "Error removing user from room")
	}

	return nil
}

func (ph *ProductionHandler) HandleGetUserEditModal(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: GET edit modal unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	user := ph.UserList.GetUserById(userUUID)
	if user == nil {
		log.Printf("[ERROR]: user with id %v not found in userList", userUUID)
		return c.String(http.StatusInternalServerError, "User not found")
	}

	data := struct {
		EditModal  bool
		Name       string
		AvatarPath string
	}{
		EditModal:  true,
		Name:       user.Name,
		AvatarPath: user.AvatarPath,
	}

	return c.Render(http.StatusOK, "user-edit-modal", data)
}

func (ph *ProductionHandler) HandlePostUserEdit(c echo.Context) error {
	userUUID, err := utils.GetAndValidateCookieJWT(c)
	if err != nil {
		log.Printf("[ERROR]: GET edit modal unauthorized request. ErrMsg: %v", err)
		notificationData := models.Notification{
			IsError: true,
			Content: "session token no longer valid",
		}
		c.Render(http.StatusUnprocessableEntity, "notification", notificationData)
		return c.Render(http.StatusUnauthorized, "login-oob", nil)
	}

	userName := c.FormValue("name-input")
	if userName == "" {
		userName = "Anonymous"
	}

	file, err := c.FormFile("avatar-input")
	fileName := "static/images/avatar_placeholder.png"
	if err == nil {
		fileName, err = utils.UploadFile(file)
		if err != nil {
			log.Printf("[ERROR]: error with file upload. ErrMsg: %v", err)
			return c.String(http.StatusInternalServerError, "Error with file upload: "+err.Error())
		}
	}

	oldUser := ph.UserList.GetUserById(userUUID)
	err = os.Remove(fmt.Sprintf("%v", oldUser.AvatarPath))
	if err != nil {
		log.Printf("[ERROR]: could not remove image %v. ErrMsg: %v", oldUser.AvatarPath, err)
	}

	user := ph.UserList.EditUser(userUUID, userName, fileName)
	data := struct {
		Name       string
		AvatarPath string
	}{
		Name:       user.Name,
		AvatarPath: user.AvatarPath,
	}

	log.Printf("[INFO]: Changed username to %v and avatarpath to %v for user: %v", user.Name, user.AvatarPath, user.Id)
	notificationData := models.Notification{
		IsInfo:  true,
		Content: "username and/or avatar edited",
	}
	c.Render(http.StatusOK, "notification", notificationData)
	// send update to all rooms
	return c.Render(http.StatusOK, "user-info", data)
}

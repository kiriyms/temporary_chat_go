package handlers

import (
	"fmt"
	"net/http"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/Kirill-Sirotkin/temporary_chat_go/utils"
	"github.com/labstack/echo/v4"
)

type ProductionHandler struct {
	UserList *models.UserList
	RoomList *models.RoomList
}

func NewProductionHandler(ul *models.UserList, rl *models.RoomList) handler {
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
		// handle error: no user in memory list
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
	room := models.NewRoom(userUUID)
	ph.RoomList.AddRoom(room)

	token, err := utils.CreateJWT(userUUID)
	if err != nil {
		fmt.Println("Could not create JWT!!!")
		return c.String(http.StatusInternalServerError, "Could not create JWT!")
	}

	jwtCookie := new(http.Cookie)
	jwtCookie.Name = "jwt"
	jwtCookie.Value = token
	c.SetCookie(jwtCookie)

	data := struct {
		UserRooms []*models.Room
	}{
		ph.RoomList.GetUserRooms(userUUID),
	}

	return c.Render(http.StatusOK, "room-list", data)
}

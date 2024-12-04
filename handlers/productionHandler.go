package handlers

import (
	"fmt"
	"net/http"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/Kirill-Sirotkin/temporary_chat_go/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductionHandler struct {
	UserList *models.UserList
}

func NewProductionHandler(ul *models.UserList) handler {
	return &ProductionHandler{
		UserList: ul,
	}
}

func (ph *ProductionHandler) HandleGetMain(c echo.Context) error {
	// make proper Data struct in models directory
	Data := models.NewIndexData()

	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		fmt.Println(err)
		return c.Render(http.StatusOK, "index", Data)
	}

	userIdStr, err := utils.ValidateJWT(jwtCookie.Value)
	if err != nil {
		// jwt error (expired, etc.)
		fmt.Println(err)
		return c.Render(http.StatusOK, "index", Data)
	}

	userUUID, err := uuid.Parse(userIdStr)
	if err != nil {
		// uuid format is wrong
		fmt.Println(err)
		return c.Render(http.StatusOK, "index", Data)
	}

	Data.Rooms = true
	user := ph.UserList.GetUserById(userUUID)
	Data.UserData = user

	return c.Render(http.StatusOK, "index", Data)
}

func (ph *ProductionHandler) HandleGetRooms(c echo.Context) error {
	return c.Render(http.StatusOK, "rooms", nil)
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

	return c.Render(http.StatusOK, "rooms", newUser)
}

// On POST /room renew jwt and add +10 mins to Expiration

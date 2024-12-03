package handlers

import (
	"fmt"
	"net/http"

	"github.com/Kirill-Sirotkin/temporary_chat_go/utils"
	"github.com/labstack/echo/v4"
)

type ProductionHandler struct{}

func NewProductionHandler() handler {
	return &ProductionHandler{}
}

func (ph *ProductionHandler) HandleGetMain(c echo.Context) error {
	// make proper Data struct somewhere
	Data := struct {
		Rooms bool
	}{false}

	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		// no cookie
		return c.Render(http.StatusOK, "index", Data)
	}

	// validate cookie sessionID, return Data with error (e.g. session expired)
	fmt.Println(sessionId.Value)

	// if valid:
	// Data.Rooms = true

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

	return c.Render(http.StatusOK, "spinner-svg", nil)
}

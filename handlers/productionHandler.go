package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/labstack/echo/v4"
)

type ProductionHandler struct{}

func NewProductionHandler() handler {
	return &ProductionHandler{}
}

func (ph *ProductionHandler) HandleGetMainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (ph *ProductionHandler) HandleGetRoomsPage(c echo.Context) error {
	return c.String(http.StatusOK, "Rooms Page!")
}

func (ph *ProductionHandler) HandlePostProfile(c echo.Context) error {
	c.Request().ParseMultipartForm(10 << 20)

	fmt.Println(c.Request().FormValue("name-input"))

	file, handler, err := c.Request().FormFile("avatar-input")
	if err != nil {
		if err == http.ErrMissingFile {
			return c.String(http.StatusBadRequest, "No file submitted!")
		} else {
			return c.String(http.StatusInternalServerError, "Error retrieving the file")
		}
	}

	nsec := time.Now().UnixNano()
	filename := handler.Filename[:len(handler.Filename)-len(filepath.Ext(handler.Filename))] +
		strconv.Itoa(int(nsec)) +
		filepath.Ext(handler.Filename)
	filePath := filepath.Join("uploads", filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error saving the file")
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		return c.String(http.StatusInternalServerError, "Error saving the file")
	}

	return c.String(http.StatusOK, "Profle Posted!")
}

func (ph *ProductionHandler) HandlePutProfile(c echo.Context) error {
	return c.String(http.StatusOK, "Profile Updated!!")
}

func (ph *ProductionHandler) HandlePostRoom(c echo.Context) error {
	return c.String(http.StatusOK, "Room Created!")
}

func (ph *ProductionHandler) HandlePostChatMessage(c echo.Context) error {
	return c.String(http.StatusOK, "Message Posted!")
}

func (ph *ProductionHandler) HandleGetWebSocketConn(c echo.Context) error {
	return c.String(http.StatusOK, "Websocket Established!")
}

func (ph *ProductionHandler) HandleGetRoomChat(c echo.Context) error {
	return c.String(http.StatusOK, "Room Chat!")
}

func (ph *ProductionHandler) HandlePostToken(c echo.Context) error {
	token := &models.Token{}
	err := json.NewDecoder(c.Request().Body).Decode(token)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad token payload. Expected JSON object with application/json Content-Type")
	}
	if token.Token == "" {
		return c.Render(http.StatusOK, "login", nil)
	}

	return c.Render(http.StatusOK, "rooms", nil)
}

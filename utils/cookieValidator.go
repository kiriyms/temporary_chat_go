package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetAndValidateCookieJWT(c echo.Context) (uuid.UUID, error) {
	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		fmt.Println(err)
		return uuid.Nil, err
	}

	userIdStr, err := ValidateJWT(jwtCookie.Value)
	if err != nil {
		// jwt error (expired, etc.)
		fmt.Println(err)
		return uuid.Nil, err
	}

	return userIdStr, nil
}

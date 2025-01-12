package utils

import (
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetAndValidateCookieJWT(c echo.Context) (uuid.UUID, error) {
	jwtCookie, err := c.Cookie("jwt")
	if err != nil {
		return uuid.Nil, err
	}

	userIdStr, err := ValidateJWT(jwtCookie.Value)
	log.Printf("[INFO]: submitted jwt: %v", jwtCookie.Value)
	if err != nil {
		return uuid.Nil, err
	}

	return userIdStr, nil
}

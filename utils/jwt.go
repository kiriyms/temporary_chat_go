package utils

import (
	"errors"
	"os"
	"time"

	"github.com/Kirill-Sirotkin/temporary_chat_go/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJWT(uId uuid.UUID) (string, error) {
	claims := models.JwtClaims{
		UserId: uId.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(180 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenStr, err
}

func ValidateJWT(t string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(t, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*models.JwtClaims)
	if !ok {
		return uuid.Nil, errors.New("unknown claims type")
	}

	userUUID, err := uuid.Parse(claims.UserId)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

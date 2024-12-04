package models

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	UserId string `json:"UserId"`
	jwt.RegisteredClaims
}

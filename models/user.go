package models

import "github.com/google/uuid"

type User struct {
	Id         uuid.UUID
	Name       string
	AvatarPath string
}

func NewUser(name, avatarPath string) *User {
	return &User{
		Id:         uuid.New(),
		Name:       name,
		AvatarPath: avatarPath,
	}
}

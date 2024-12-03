package models

import (
	"github.com/google/uuid"
)

type Room struct {
	Id       uuid.UUID
	Messages []Message
	Users    []uuid.UUID
}

func NewRoom(creatorId uuid.UUID) *Room {
	room := &Room{
		Id:       uuid.New(),
		Messages: make([]Message, 0),
		Users:    make([]uuid.UUID, 0),
	}
	room.Users = append(room.Users, creatorId)

	return room
}

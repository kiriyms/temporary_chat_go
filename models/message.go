package models

import "github.com/google/uuid"

type Message struct {
	SenderId uuid.UUID
	Content  string
}

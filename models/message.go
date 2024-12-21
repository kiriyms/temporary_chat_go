package models

import "github.com/google/uuid"

type Message struct {
	Id      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

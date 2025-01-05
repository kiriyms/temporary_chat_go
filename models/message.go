package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

type MessageWithTimeCode struct {
	Id       uuid.UUID
	Content  string
	TimeCode time.Time
}

func NewMessageWithTimeCode(id uuid.UUID, content string) *MessageWithTimeCode {
	return &MessageWithTimeCode{
		Id:       id,
		Content:  content,
		TimeCode: time.Now(),
	}
}

package models

import (
	"github.com/google/uuid"
)

type Message struct {
	ChatID       uuid.UUID `json:"chat_id"`
	SenderID     uuid.UUID `json:"sender_id"`
	MsgID        uuid.UUID `json:"msg_id"`
	Payload      string    `json:"payload"`
	CreatedAt    int64     `json:"created_at"`
	SessionStart bool      `json:"session_start,omitempty"`
}

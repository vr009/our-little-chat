package models

import (
	"github.com/google/uuid"
)

type Message struct {
	ChatID    uuid.UUID `json:"chat_id" bson:"chat_id"`
	MsgID     uuid.UUID `json:"msg_id,omitempty" bson:"msg_id"`
	SenderID  uuid.UUID `json:"sender_id" bson:"sender_id"`
	Payload   string    `json:"payload" bson:"payload"`
	CreatedAt int64     `json:"created_at,omitempty" bson:"created_at"`
}

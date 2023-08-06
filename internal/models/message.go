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

type Messages []Message

func (m Messages) Len() int           { return len(m) }
func (m Messages) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Messages) Less(i, j int) bool { return m[i].CreatedAt < m[j].CreatedAt }

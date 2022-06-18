package models

import (
	"github.com/google/uuid"
)

type Message struct {
	ChatID     uuid.UUID `json:"chatID"`
	MsgID      uuid.UUID `json:"-"`
	SenderID   uuid.UUID `json:"senderID"`
	ReceiverID uuid.UUID `json:"receiverID"`
	Payload    string    `json:"payload"`
	CreatedAt  float64   `json:"-"`
}

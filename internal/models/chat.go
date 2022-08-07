package models

import (
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ChatID      uuid.UUID `json:"chat_id" bson:"chat_id"`
	Owner       uuid.UUID `json:"owner" bson:"owner"`
	Opponent    uuid.UUID `json:"opponent" bson:"opponent"`
	LastMessage time.Time `json:"last_message" bson:"last_message"`
}

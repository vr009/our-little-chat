package models

import (
	"github.com/google/uuid"
)

type Chat struct {
	ChatID       uuid.UUID   `json:"chat_id,omitempty"`
	Participants []uuid.UUID `json:"participants,omitempty"`
	PhotoURL     string      `json:"photo_url,omitempty"`
	CreatedAt    int64       `json:"created_at,omitempty"`
	LastMessage  Message     `json:"last_message,omitempty"`
}

package models

import (
	"github.com/google/uuid"
)

type Chat struct {
	ChatID       uuid.UUID   `json:"chat_id,omitempty"`
	Participants []uuid.UUID `json:"participants,omitempty"`
}

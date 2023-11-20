package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
}

const (
	ActivationSession = "activation-session"
	PlainSession      = "plain-session"
)

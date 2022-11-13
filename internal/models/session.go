package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time
	ExpireAt  time.Time
}

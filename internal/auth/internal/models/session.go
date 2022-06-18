package models

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

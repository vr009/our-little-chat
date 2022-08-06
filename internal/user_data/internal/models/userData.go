package models

import (
	"github.com/google/uuid"
	"time"
)

type UserData struct {
	UserID uuid.UUID `json:"user_id"`

	Nickname string `json:"nickname"`

	LastAuth time.Time `json:"last_auth"`

	Registered time.Time `json:"registered"`

	Avatar []string `json:"avatar"`

	ContactList []string `json:"contact_list"`
}

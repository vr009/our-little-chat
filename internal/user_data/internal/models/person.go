package models

import (
	"github.com/google/uuid"
	"time"
)

type Person struct {
	UserId      uuid.UUID `json:"user_id,omitempty"`
	Nickname    string    `json:"nickname"`
	LastAuth    time.Time `json:"last_auth"`
	Registered  time.Time `json:"registered"`
	Avatar      []string  `json:"avatar"`
	ContactList []string  `json:"contact_list"`
}

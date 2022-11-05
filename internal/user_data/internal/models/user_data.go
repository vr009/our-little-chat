package models

import (
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	UserID      uuid.UUID `json:"user_id,omitempty"`
	Nickname    string    `json:"nickname"`
	Password    string    `json:"password"`
	LastAuth    time.Time `json:"last_auth,omitempty"`
	Registered  time.Time `json:"registered,omitempty"`
	Avatar      []string  `json:"avatar"`
	ContactList []string  `json:"contact_list"`
}

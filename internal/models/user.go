package models

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Nickname string    `json:"nickname,omitempty"`
	Name     string    `json:"name,omitempty"`
	Surname  string    `json:"surname,omitempty"`
	Password string    `json:"password,omitempty"`
}

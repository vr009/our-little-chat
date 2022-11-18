package models

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `json:"user_id,omitempty" bson:"uuid"`
	Nickname string    `json:"nickname,omitempty"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

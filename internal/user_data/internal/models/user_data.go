package models

import (
	"time"

	"our-little-chatik/internal/models"
)

type UserData struct {
	models.User
	LastAuth   time.Time `json:"last_auth,omitempty"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
}

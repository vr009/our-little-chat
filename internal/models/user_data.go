package models

import (
	"time"
)

type UserData struct {
	User
	LastAuth   time.Time `json:"last_auth,omitempty"`
	Registered time.Time `json:"registered,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
}

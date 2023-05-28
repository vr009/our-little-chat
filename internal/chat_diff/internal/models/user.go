package models

import (
	"our-little-chatik/internal/models"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID `json:"user_id"`
}

type ChatDiffUser struct {
	User    User
	Updates chan []models.ChatItem
}

package models

import (
	"our-little-chatik/internal/models"
)

type ChatDiffUser struct {
	User    models.User
	Updates chan []models.ChatItem
}

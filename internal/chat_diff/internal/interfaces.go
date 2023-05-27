package internal

import (
	models2 "our-little-chatik/internal/chat_diff/internal/models"
	"our-little-chatik/internal/models"
)

type ChatDiffRepo interface {
	FetchUpdates(user models.User) []models.ChatItem
}

type Manager interface {
	Work()
	AddChatUser(user *models2.ChatDiffUser) *models2.ChatDiffUser
}

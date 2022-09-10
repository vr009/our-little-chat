package internal

import (
	"github.com/google/uuid"
	models2 "our-little-chatik/internal/chat_diff/internal/models"
)

type ChatDiffUsecase interface {
	GetUpdates(user models2.ChatUser) []models2.ChatItem
}

type ChatDiffRepo interface {
	FetchUpdates(user models2.ChatUser) []models2.ChatItem
}

type Manager interface {
	Work()
	AddChatUser(user *models2.ChatUser) *models2.ChatUser
}

type TokenResolver interface {
	ResolveToken(token string) (uuid.UUID, error)
}

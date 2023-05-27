package internal

import (
	"our-little-chatik/internal/models"

	"github.com/google/uuid"
)

type ChatDiffUsecase interface {
	GetUpdates(user models.User) []models.ChatItem
}

type ChatDiffRepo interface {
	FetchUpdates(user models.User) []models.ChatItem
}

type Manager interface {
	Work()
	AddChatUser(user *models.User) *models.User
}

type TokenResolver interface {
	ResolveToken(token string) (uuid.UUID, error)
}

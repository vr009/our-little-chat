package internal

import (
	"context"
	"our-little-chatik/internal/models"
)

type DiffRepo interface {
	SubscribeToChats(ctx context.Context, chats []models.Chat) (chan models.Message, error)
}

package internal

import (
	"context"
	"our-little-chatik/internal/models"
)

type PeerRepo interface {
	SendToChannel(ctx context.Context, msg models.Message, chatChannel string)
	CheckUserExists(ctx context.Context, user string, userSet string) (bool, error)
	CreateUser(ctx context.Context, user string, userSet string) error
	RemoveUser(ctx context.Context, user string, userSet string)
	StartSubscriber(ctx context.Context, messageChan chan models.Message, chatChannel string)
	SaveMessage(message models.Message) error
}

type DiffRepo interface {
	SubscribeToChats(ctx context.Context, chats []models.Chat) (chan models.Message, error)
}

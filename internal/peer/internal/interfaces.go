package internal

import (
	"context"
	"our-little-chatik/internal/models"
)

type PeerRepo interface {
	CheckUserExists(ctx context.Context, user string, userSet string) (bool, error)
	CreateUser(ctx context.Context, user string, userSet string) error
	RemoveUser(ctx context.Context, user string, userSet string)
	SaveMessage(message models.Message) error
}

type MessageBus interface {
	SubscribeOnChatMessages(ctx context.Context, chatChannel string, readyChan chan struct{}) chan models.Message
	SendMessageToChannel(ctx context.Context, msg models.Message, chatChannel string)
}

type DiffRepo interface {
	SubscribeToChats(ctx context.Context, chats []models.Chat) (chan models.Message, error)
}

package internal

import "our-little-chatik/internal/models"

type QueueRepo interface {
	FetchAllMessages() ([]models.Message, error)
	FetchChatListUpdate() ([]models.ChatItem, error)
}

type PersistantRepo interface {
	PersistAllMessages(msgs []models.Message) error
	PersistChatListUpdate(chats []models.ChatItem) error
}

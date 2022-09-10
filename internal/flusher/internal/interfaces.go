package internal

import "our-little-chatik/internal/models"

type QueueRepo interface {
	FetchAllMessages() ([]models.Message, error)
	FetchAllChats() ([]models.Chat, error)
}

type PersistantRepo interface {
	PersistAllMessages(msgs []models.Message) error
	PersistAllChats(chats []models.Chat) error
}

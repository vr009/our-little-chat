package internal

import "our-little-chatik/internal/models"

type QueueRepo interface {
	FetchAll() ([]models.Message, error)
	FetchAllChats() ([]models.Chat, error)
}

type PersistantRepo interface {
	PersistAll(msgs []models.Message) error
	PersistAllChats(chats []models.Chat) error
}

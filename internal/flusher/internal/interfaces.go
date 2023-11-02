package internal

import "our-little-chatik/internal/models"

type QueueRepo interface {
	FetchAllMessages() ([]models.Message, error)
	FetchAllLastMessagesOfChats() ([]models.Message, error)
}

type PersistantRepo interface {
	PersistAllMessages(msgs []models.Message) error
	PersistAllLastChatMessages(msgs []models.Message) error
}

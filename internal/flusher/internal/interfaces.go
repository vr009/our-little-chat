package internal

import "our-little-chatik/internal/models"

type QueueRepo interface {
	FetchAllMessages() ([]models.Message, error)
}

type PersistantRepo interface {
	PersistAllMessages(msgs []models.Message) error
}

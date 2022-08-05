package internal

import "our-little-chatik/internal/models"

type QueueRepo interface {
	FetchAll() ([]models.Message, error)
}

type PersistantRepo interface {
	PersistAll(msgs []models.Message) error
}

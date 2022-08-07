package internal

import (
	"our-little-chatik/internal/models"
)

type ChatRepo interface {
	InsertMessages(mes []models.Message) error
	GetChat(chat models.Chat, opts models.Opts) ([]models.Message, error)
}

type QueueRepo interface {
	GetFreshChat(chat models.Chat) ([]models.Message, error)
}

type ChatUseCase interface {
	SaveMessages(msgs []models.Message) error
	FetchChat(chat models.Chat, opts models.Opts) ([]models.Message, error)
}

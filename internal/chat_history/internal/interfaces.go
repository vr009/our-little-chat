package internal

import (
	"our-little-chatik/internal/models"
)

type ChatRepo interface {
	GetChat(chat models.Chat, opts models.Opts) ([]models.Message, error)
}

type QueueRepo interface {
	GetFreshChat(chat models.Chat) ([]models.Message, error)
}

type ChatUseCase interface {
	FetchChat(chat models.Chat, opts models.Opts) ([]models.Message, error)
}

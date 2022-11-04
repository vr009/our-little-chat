package internal

import (
	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"
)

type ChatRepo interface {
	GetChatMessages(chat models2.Chat, opts models.Opts) ([]models.Message, error)
	FetchChatList(user models.User) ([]models.ChatItem, error)
}

type QueueRepo interface {
	InsertChat(models2.Chat) (models2.Chat, error)
	GetFreshMessagesFromChat(chat models2.Chat) ([]models.Message, error)
}

type ChatUseCase interface {
	CreateNewChat(chat models2.Chat) (models2.Chat, error)
	ActivateChat(chat models2.Chat) error
	FetchChatMessages(chat models2.Chat, opts models.Opts) ([]models.Message, error)
	GetChatList(user models.User) ([]models.ChatItem, error)
}

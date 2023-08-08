package internal

import (
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
)

type ChatRepo interface {
	GetChatMessages(chat models.Chat, opts models.Opts) (models.Messages, error)
	FetchChatList(user models.User) ([]models.ChatItem, error)
	InsertChat(models.Chat) error
	UpdateChat(chat models.Chat, updateOpts models2.UpdateOptions) error
	GetChat(chat models.Chat) (models.Chat, error)
}

type QueueRepo interface {
	GetFreshMessagesFromChat(chat models.Chat) (models.Messages, error)
}

type ChatUseCase interface {
	CreateNewChat(chat models.Chat) (models.Chat, error)
	FetchChatMessages(chat models.Chat, opts models.Opts) (models.Messages, error)
	GetChatList(user models.User) ([]models.ChatItem, error)
	UpdateChat(chat models.Chat, updateOpts models2.UpdateOptions) error
	GetChat(chat models.Chat) (models.Chat, error)
}

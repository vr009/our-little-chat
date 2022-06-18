package internal

import (
	"our-little-chatik/internal/peer/internal/models"
)

type WebSocketWorker interface {
	Read()
	Write()
	Close()
}

type PeerUsecase interface {
	SendMessage(msg *models.Message, chat *models.Chat) error
	FetchMessages(chat *models.Chat) ([]models.Message, error)
}

type PeerRepo interface {
	SendPayload(msg *models.Message, chat *models.Chat) error
	FetchUpdates(chat *models.Chat) ([]models.Message, error)
}

type MessageManager interface {
	Work()
	// EnqueueChatIfNotExists enqueues a passed Chat to an internal queue of chats.
	// If the chat already exists it finds it and return.
	EnqueueChatIfNotExists(chat *models.Chat) *models.Chat
	// DequeueChat dequeues chat from internal common queue
	DequeueChat(chat *models.Chat)
	// EnqueueChat enqueues chat in internal common queue
	EnqueueChat(chat *models.Chat)
}

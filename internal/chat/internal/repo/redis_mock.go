package repo

import (
	"fmt"
	"github.com/google/uuid"
	"our-little-chatik/internal/models"
)

type RedisMock struct {
	Msgs models.Messages
	ID   uuid.UUID
}

func (m RedisMock) GetFreshMessagesFromChat(chat models.Chat) (models.Messages, error) {
	if chat.ChatID != m.ID {
		return nil, fmt.Errorf("not found")
	}
	return m.Msgs, nil
}

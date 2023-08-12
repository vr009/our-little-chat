package repo

import (
	"fmt"
	"github.com/google/uuid"
	"our-little-chatik/internal/models"
	"sort"
)

type RedisMock struct {
	Msgs models.Messages
	ID   uuid.UUID
}

func (m RedisMock) GetChatMessages(chat models.Chat,
	opts models.Opts) (models.Messages, error) {
	if chat.ChatID != m.ID {
		return nil, fmt.Errorf("not found")
	}
	sort.Sort(m.Msgs)
	return m.Msgs, nil
}

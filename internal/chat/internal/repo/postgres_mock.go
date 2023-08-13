package repo

import (
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"sort"
)

type PostgresMock struct {
	Msgs     models.Messages
	ChatList []models.ChatItem
	Chat     models.Chat
}

func (m PostgresMock) GetChatMessages(Chat models.Chat, opts models.Opts) (models.Messages, error) {
	sort.Sort(m.Msgs)
	return m.Msgs, nil
}

func (m PostgresMock) FetchChatList(user models.User) ([]models.ChatItem, error) {
	return m.ChatList, nil
}

func (m PostgresMock) CreateChat(Chat models.Chat) error {
	return nil
}

func (m PostgresMock) UpdateChat(Chat models.Chat, updateOpts models2.UpdateOptions) error {
	return nil
}

func (m PostgresMock) GetChat(Chat models.Chat) (models.Chat, error) {
	return m.Chat, nil
}

func (m PostgresMock) DeleteMessage(message models.Message) error {
	return nil
}

func (m PostgresMock) DeleteChat(Chat models.Chat) error {
	return nil
}

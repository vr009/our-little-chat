package usecase

import (
	"fmt"

	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"

	"github.com/google/uuid"
)

type ChatUseCase struct {
	repo  internal.ChatRepo
	queue internal.QueueRepo
}

func NewChatUseCase(rep internal.ChatRepo, queue internal.QueueRepo) *ChatUseCase {
	return &ChatUseCase{repo: rep, queue: queue}
}

func (ch *ChatUseCase) FetchChatMessages(chat models2.Chat, opts models.Opts) ([]models.Message, error) {
	msgs, err := ch.queue.GetFreshMessagesFromChat(chat)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	oldMsgs, err := ch.repo.GetChatMessages(chat, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	msgs = append(msgs, oldMsgs...)
	return msgs, nil
}

func (ch *ChatUseCase) GetChatList(user models.User) ([]models.ChatItem, error) {
	return ch.repo.FetchChatList(user)
}

func (ch *ChatUseCase) CreateNewChat(chat models2.Chat) (models2.Chat, error) {
	chat.ChatID = uuid.New()
	return ch.queue.InsertChat(chat)
}

func (ch *ChatUseCase) ActivateChat(chat models2.Chat) error {
	_, err := ch.queue.InsertChat(chat)
	return err
}

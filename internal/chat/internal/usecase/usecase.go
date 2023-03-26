package usecase

import (
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"

	"github.com/golang/glog"
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
		glog.Error(err)
	}
	oldMsgs, err := ch.repo.GetChatMessages(chat, opts)
	if err != nil {
		glog.Error(err)
	}
	msgs = append(msgs, oldMsgs...)
	return msgs, nil
}

func (ch *ChatUseCase) GetChatList(user models.User) ([]models.ChatItem, error) {
	return ch.repo.FetchChatList(user)
}

func (ch *ChatUseCase) CreateNewChat(chat models2.Chat) (models2.Chat, error) {
	chat.ChatID = uuid.New()
	err := ch.repo.InsertChat(chat)
	if err != nil {
		glog.Error(err.Error())
	}
	return ch.queue.InsertChat(chat)
}

func (ch *ChatUseCase) ActivateChat(chat models2.Chat) error {
	_, err := ch.queue.InsertChat(chat)
	return err
}

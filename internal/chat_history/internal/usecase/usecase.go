package usecase

import (
	"fmt"
	"our-little-chatik/internal/chat_history/internal"
	"our-little-chatik/internal/models"
)

type ChatUseCase struct {
	repo  internal.ChatRepo
	queue internal.QueueRepo
}

func NewChatUseCase(rep internal.ChatRepo, queue internal.QueueRepo) *ChatUseCase {
	return &ChatUseCase{repo: rep, queue: queue}
}

func (ch *ChatUseCase) FetchChat(chat models.Chat, opts models.Opts) ([]models.Message, error) {
	msgs, err := ch.queue.GetFreshChat(chat)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	oldMsgs, err := ch.repo.GetChat(chat, opts)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	msgs = append(msgs, oldMsgs...)
	return msgs, nil
}

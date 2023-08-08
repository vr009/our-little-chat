package usecase

import (
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"sort"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type ChatUseCase struct {
	repo  internal.ChatRepo
	queue internal.QueueRepo
}

func NewChatUseCase(rep internal.ChatRepo, queue internal.QueueRepo) *ChatUseCase {
	return &ChatUseCase{repo: rep, queue: queue}
}

func (ch *ChatUseCase) FetchChatMessages(chat models.Chat,
	opts models.Opts) (models.Messages, error) {
	msgs, err := ch.queue.GetFreshMessagesFromChat(chat)
	if err != nil {
		slog.Error(err.Error())
	}
	oldMsgs, err := ch.repo.GetChatMessages(chat, opts)
	if err != nil {
		slog.Error(err.Error())
	}
	msgs = append(msgs, oldMsgs...)
	sort.Sort(msgs)
	return msgs, nil
}

func (ch *ChatUseCase) GetChatList(user models.User) ([]models.ChatItem, error) {
	return ch.repo.FetchChatList(user)
}

func (ch *ChatUseCase) CreateNewChat(chat models.Chat) (models.Chat, error) {
	chat.ChatID = uuid.New()
	chat.CreatedAt = time.Now().Unix()
	err := ch.repo.InsertChat(chat)
	if err != nil {
		slog.Error(err.Error())
	}
	return chat, err
}

func (ch *ChatUseCase) UpdateChat(chat models.Chat, updateOpts models2.UpdateOptions) error {
	return ch.repo.UpdateChat(chat, updateOpts)
}

func (ch *ChatUseCase) GetChat(chat models.Chat) (models.Chat, error) {
	return ch.repo.GetChat(chat)
}

package usecase

import (
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
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

func (ch *ChatUseCase) GetChatMessages(chat models.Chat,
	opts models.Opts) (models.Messages, error) {
	msgs, err := ch.queue.GetChatMessages(chat, opts)
	if err != nil {
		slog.Error(err.Error())
	}
	if len(msgs) < int(opts.Limit) {
		opts.Limit = opts.Limit - int64(len(msgs))
		oldMsgs, err := ch.repo.GetChatMessages(chat, opts)
		if err != nil {
			slog.Error(err.Error())
		}
		msgs = append(msgs, oldMsgs...)
	}
	return msgs, nil
}

func (ch *ChatUseCase) GetChatList(user models.User) ([]models.ChatItem, error) {
	return ch.repo.FetchChatList(user)
}

func (ch *ChatUseCase) CreateChat(chat models.Chat) (models.Chat, error) {
	chat.ChatID = uuid.New()
	chat.CreatedAt = time.Now().Unix()
	err := ch.repo.CreateChat(chat)
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

func (ch *ChatUseCase) DeleteChat(chat models.Chat) error {
	return ch.repo.DeleteChat(chat)
}

func (ch *ChatUseCase) DeleteMessage(message models.Message) error {
	return ch.repo.DeleteMessage(message)
}

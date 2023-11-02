package internal

import (
	"context"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
)

type ChatRepo interface {
	GetChatMessages(ctx context.Context, chat models.Chat, opts models.Opts) (models.Messages, models.StatusCode)
	FetchChatList(ctx context.Context, user models.User) ([]models.Chat, models.StatusCode)
	CreateChat(ctx context.Context, chat models.Chat, chatNames map[string]string) models.StatusCode
	GetChat(ctx context.Context, chat models.Chat) (models.Chat, models.StatusCode)
	DeleteMessage(ctx context.Context, message models.Message) models.StatusCode
	DeleteChat(ctx context.Context, chat models.Chat) models.StatusCode
	RemoveUserFromChat(ctx context.Context,
		chat models.Chat, users ...models.User) models.StatusCode
	AddUsersToChat(ctx context.Context,
		chat models.Chat, chatNames map[string]string, users ...models.User) models.StatusCode
	UpdateChatPhotoURL(ctx context.Context, chat models.Chat,
		photoURL string) models.StatusCode
}

type QueueRepo interface {
	GetChatMessages(chat models.Chat, opts models.Opts) (models.Messages, models.StatusCode)
	GetChatsLastMessages(chatList []models.Chat) (models.Messages, models.StatusCode)
}

type ChatUseCase interface {
	CreateChat(ctx context.Context, chat models2.CreateChatRequest) (models.Chat, models.StatusCode)
	GetChatMessages(ctx context.Context, chat models.Chat, opts models.Opts) (models.Messages, models.StatusCode)
	GetChatList(ctx context.Context, user models.User) ([]models.Chat, models.StatusCode)
	GetChat(ctx context.Context, chat models.Chat) (models.Chat, models.StatusCode)
	DeleteChat(ctx context.Context, chat models.Chat) models.StatusCode
	DeleteMessage(ctx context.Context, message models.Message) models.StatusCode
	RemoveUserFromChat(ctx context.Context,
		chat models.Chat, users ...models.User) models.StatusCode
	AddUsersToChat(ctx context.Context,
		chat models.Chat, users ...models.User) models.StatusCode
	UpdateChatPhotoURL(ctx context.Context, chat models.Chat,
		photoURL string) models.StatusCode
}

type UserDataInteractor interface {
	GetUser(user models.User) (models.User, models.StatusCode)
}

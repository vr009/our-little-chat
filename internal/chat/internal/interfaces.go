package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
)

type ChatRepo interface {
	GetChatMessages(chat models.Chat, opts models.Opts) (models.Messages, error)
	FetchChatList(user models.User) ([]models.ChatItem, error)
	CreateChat(models.Chat) error
	UpdateChat(chat models.Chat, updateOpts models2.UpdateOptions) error
	GetChat(chat models.Chat) (models.Chat, error)
	DeleteMessage(message models.Message) error
	DeleteChat(chat models.Chat) error
}

type QueueRepo interface {
	GetChatMessages(chat models.Chat, opts models.Opts) (models.Messages, error)
}

type ChatUseCase interface {
	CreateChat(chat models.Chat) (models.Chat, error)
	GetChatMessages(chat models.Chat, opts models.Opts) (models.Messages, error)
	GetChatList(user models.User) ([]models.ChatItem, error)
	UpdateChat(chat models.Chat, updateOpts models2.UpdateOptions) error
	GetChat(chat models.Chat) (models.Chat, error)
	DeleteChat(chat models.Chat) error
	DeleteMessage(message models.Message) error
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

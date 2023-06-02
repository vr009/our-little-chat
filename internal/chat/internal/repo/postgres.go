package repo

import (
	"context"

	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

const (
	InsertChatParticipantsQuery = "INSERT INTO chat_participants VALUES ($1, $2)"
	GetMessagesQuery            = "SELECT msg_id, sender_id, payload, created_at  FROM messages WHERE chat_id=$1 ORDER BY created_at ASC OFFSET $2 LIMIT $3"
	FetchChatListQuery          = "SELECT chat_id FROM chat_participants WHERE participant_id=$1" //TODO complete it
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

func (pr PostgresRepo) GetChatMessages(chat models2.Chat, opts models.Opts) ([]models.Message, error) {
	ctx := context.Background()
	rows, err := pr.pool.Query(ctx, GetMessagesQuery, chat.ChatID, opts.Page, opts.Limit)
	if err != nil {
		return nil, err
	}

	msgs := make([]models.Message, 0)
	for rows.Next() {
		msg := models.Message{}
		err := rows.Scan(&msg.ChatID, &msg.SenderID, &msg.Payload, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (pr PostgresRepo) FetchChatList(user models.User) ([]models.ChatItem, error) {
	ctx := context.Background()
	rows, err := pr.pool.Query(ctx, FetchChatListQuery, user.UserID)
	if err != nil {
		return nil, err
	}

	chatList := make([]models.ChatItem, 0)
	for rows.Next() {
		chat := models.ChatItem{}
		err := rows.Scan(&chat.ChatID)
		if err != nil {
			return nil, err
		}

		chatList = append(chatList, chat)
	}

	return chatList, nil
}

func (pr PostgresRepo) InsertChat(chat models2.Chat) error {
	ctx := context.Background()
	batch := &pgx.Batch{}
	for _, participant := range chat.Participants {
		batch.Queue(InsertChatParticipantsQuery, chat.ChatID, participant).Exec(func(ct pgconn.CommandTag) error {
			return nil
		})
	}

	results := pr.pool.SendBatch(ctx, batch)
	defer results.Close()
	for _, participant := range chat.Participants {
		_, err := results.Exec()
		if err != nil {
			slog.Error("Failed to add a chat user", "user", participant.String())
			return err
		}
	}
	return nil
}

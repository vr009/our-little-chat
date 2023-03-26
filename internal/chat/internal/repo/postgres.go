package repo

import (
	"context"

	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	InsertChatParticipantsQuery = "INSERT INTO chat_participants VALUES (?, ?)"
	GetMessagesQuery            = "SELECT msg_id, sender_id, payload, created_at  FROM messages WHERE chat_id=? OFFSET ? LIMIT ? ORDER BY created_at ASC"
	FetchChatListQuery          = "SELECT chat_id FROM chat_participants WHERE participant_id=?" //TODO complete it
)

type PostgresRepo struct {
	conn *pgx.Conn
}

func NewPostgresRepo(conn *pgx.Conn) *PostgresRepo {
	return &PostgresRepo{conn: conn}
}

func (pr PostgresRepo) GetChatMessages(chat models2.Chat, opts models.Opts) ([]models.Message, error) {
	ctx := context.Background()
	rows, err := pr.conn.Query(ctx, GetMessagesQuery, chat.ChatID, opts.Page, opts.Limit)
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
	rows, err := pr.conn.Query(ctx, FetchChatListQuery, user.UserID)
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

	err := pr.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}
	return nil
}

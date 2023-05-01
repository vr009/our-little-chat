package repo

import (
	"context"

	"our-little-chatik/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	InsertMsgQuery         = "INSERT INTO messages VALUES ($1, $2, $3, $4, $5)"
	UpdateChatInfo         = "UPDATE chats SET last_msg_id=$1 WHERE chat_id=$2"
	UpsertChatInfo         = "INSERT INTO chats (chat_id, last_msg_id) VALUES ($1, $2) ON CONFLICT (chat_id) DO UPDATE SET last_msg_id=EXCLUDED.last_msg_id"
	UpsertChatParticipants = "INSERT INTO chat_participants VALUES ($1, $2, $3) ON CONFLICT (chat_id, participant_id) DO UPDATE SET last_read_msg_id=EXCLUDED.last_read_msg_id"
)

type PostgresRepo struct {
	conn *pgx.Conn
}

func NewPostgresRepo(conn *pgx.Conn) *PostgresRepo {
	return &PostgresRepo{conn: conn}
}

func (pr PostgresRepo) PersistAllMessages(msgs []models.Message) error {
	ctx := context.Background()
	batch := &pgx.Batch{}
	for _, msg := range msgs {
		batch.Queue(InsertMsgQuery, msg.MsgID, msg.ChatID, msg.SenderID, msg.Payload, msg.CreatedAt).
			Exec(func(ct pgconn.CommandTag) error {
				return nil
			})
	}
	err := pr.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}

	return nil
}

func (pr PostgresRepo) PersistChatListUpdate(chats []models.ChatItem) error {
	ctx := context.Background()
	batch := &pgx.Batch{}
	for _, chat := range chats {
		batch.Queue(UpsertChatInfo, chat.ChatID, chat.MsgID).
			Exec(func(ct pgconn.CommandTag) error {
				return nil
			})
	}
	err := pr.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}
	return nil
}

func (pr PostgresRepo) PersistChatParticipants(chats []models.Chat) error {
	ctx := context.Background()
	batch := &pgx.Batch{}
	for _, chat := range chats {
		batch.Queue(UpsertChatParticipants, chat.ChatID, chat.Participant, chat.LastReadMsgID).
			Exec(func(ct pgconn.CommandTag) error {
				return nil
			})
	}
	err := pr.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}
	return nil
}

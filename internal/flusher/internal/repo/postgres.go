package repo

import (
	"context"

	"our-little-chatik/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	InsertMsgQuery = "INSERT INTO messages VALUES (?, ?, ?, ?, ?)"
	UpdateChatInfo = "UPDATE chats SET last_msg=? WHERE chat_id=?"
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
		batch.Queue(UpdateChatInfo, msg.MsgID, msg.ChatID)
	}
	err := pr.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}

	return nil
}

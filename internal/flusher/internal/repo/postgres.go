package repo

import (
	"context"

	"our-little-chatik/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	InsertMsgQuery      = "INSERT INTO messages VALUES ($1, $2, $3, $4, $5)"
	UpdateChatInfoQuery = "UPDATE chats SET last_msg_id=$1 WHERE chat_id=$2"
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

func (pr PostgresRepo) PersistAllLastChatMessages(msgs []models.Message) error {
	ctx := context.Background()
	batch := &pgx.Batch{}
	for _, msg := range msgs {
		batch.Queue(UpdateChatInfoQuery, msg.MsgID, msg.ChatID).
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

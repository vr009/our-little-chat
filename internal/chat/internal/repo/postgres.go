package repo

import (
	"context"
	"github.com/golang/glog"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

const (
	InsertChatParticipantsQuery = `INSERT INTO chat_participants VALUES ($1, $2)`
	InsertChatQuery             = `INSERT INTO chats VALUES($1, $2, $3, $4)`
	GetMessagesQuery            = `SELECT msg_id, sender_id, payload, created_at FROM messages WHERE chat_id=$1 ORDER BY created_at ASC OFFSET $2 LIMIT $3`
	GetChatInfoQuery            = `SELECT chat_id, name, photo_url, created_at FROM chats WHERE chat_id=$1`
	FetchChatListQuery          = `SELECT cp.chat_id, c.name, c.photo_url, m.payload, m.created_at FROM chat_participants AS cp 
    LEFT JOIN chats AS c ON cp.chat_id = c.chat_id 
    LEFT JOIN messages AS m ON c.last_msg_id = m.msg_id                      
                          WHERE cp.participant_id=$1`
	UpdatePhotoURLQuery     = "UPDATE chats SET photo_url=$1 WHERE chat_id=$2"
	RemoveUserFromChatQuery = "DELETE FROM chat_participants WHERE participant_id=$1 AND chat_id=$2"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

// GetChat
func (pr PostgresRepo) GetChat(chat models.Chat) (models.Chat, error) {
	ctx := context.Background()
	row := pr.pool.QueryRow(ctx, GetChatInfoQuery, chat.ChatID)
	err := row.Scan(&chat.ChatID, &chat.Name, &chat.PhotoURL, &chat.CreatedAt)
	if err != nil {
		return models.Chat{}, err
	}
	return chat, nil
}

// GetChatMessages
func (pr PostgresRepo) GetChatMessages(chat models.Chat, opts models.Opts) (models.Messages, error) {
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

// FetchChatList
func (pr PostgresRepo) FetchChatList(user models.User) ([]models.ChatItem, error) {
	ctx := context.Background()
	rows, err := pr.pool.Query(ctx, FetchChatListQuery, user.UserID)
	if err != nil {
		return nil, err
	}

	chatList := make([]models.ChatItem, 0)
	for rows.Next() {
		chat := models.ChatItem{}
		err := rows.Scan(&chat.ChatID, &chat.Name, &chat.PhotoURL,
			&chat.LastMsg, &chat.LastMessageTime)
		if err != nil {
			return nil, err
		}

		chatList = append(chatList, chat)
	}

	return chatList, nil
}

// InsertChat
func (pr PostgresRepo) InsertChat(chat models.Chat) error {
	ctx := context.Background()
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, participant := range chat.Participants {
		batch.Queue(InsertChatParticipantsQuery, chat.ChatID, participant).Exec(func(ct pgconn.CommandTag) error {
			return nil
		})
	}

	batch.Queue(InsertChatQuery, chat.ChatID, chat.Name, chat.PhotoURL, chat.CreatedAt)

	results := pr.pool.SendBatch(ctx, batch)
	defer results.Close()
	for _, participant := range chat.Participants {
		_, err := results.Exec()
		if err != nil {
			slog.Error("Failed to add a chat user", "user", participant.String())
			txErr := tx.Rollback(ctx)
			if err != nil {
				glog.Error(txErr)
			}
			return err
		}
	}
	_, err = results.Exec()
	if err != nil {
		slog.Error("Failed to add a chat")
		txErr := tx.Rollback(ctx)
		if err != nil {
			glog.Error(txErr)
		}
		return err
	}

	txErr := tx.Commit(ctx)
	if err != nil {
		glog.Error(txErr)
	}
	return nil
}

// UpdateChat
func (pr PostgresRepo) UpdateChat(chat models.Chat, updateOpts models2.UpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	switch updateOpts.Action {
	case models2.UpdatePhotoURL:
		_, err := pr.pool.Exec(ctx, UpdatePhotoURLQuery, chat.PhotoURL, chat.ChatID)
		if err != nil {
			return err
		}
	case models2.AddUsersToParticipants:
		if len(chat.Participants) > 0 {
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
		}
	case models2.RemoveUsersFromParticipants:
		if len(chat.Participants) > 0 {
			batch := &pgx.Batch{}
			for _, participant := range chat.Participants {
				batch.Queue(RemoveUserFromChatQuery, chat.ChatID, participant).Exec(func(ct pgconn.CommandTag) error {
					return nil
				})
			}
			results := pr.pool.SendBatch(ctx, batch)
			defer results.Close()
			for _, participant := range chat.Participants {
				_, err := results.Exec()
				if err != nil {
					slog.Error("Failed to remove a chat user", "user", participant.String())
					return err
				}
			}
		}
	}
	return nil
}

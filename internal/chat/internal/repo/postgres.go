package repo

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"log"
	"our-little-chatik/internal/models"
	"sort"
)

const (
	CreateChatParticipantsQuery = `INSERT INTO chat_participants VALUES ($1, $2, $3)`
	CreateChatQuery             = `INSERT INTO chats VALUES($1, $2, $3)`
	GetChatMessagesQuery        = `SELECT msg_id, sender_id, payload, created_at FROM messages WHERE chat_id=$1 ORDER BY created_at ASC OFFSET $2 LIMIT $3`
	GetChatInfoQuery            = `SELECT c.chat_id, cp.chat_name, c.photo_url, c.created_at FROM chats AS c
    LEFT JOIN chat_participants AS cp ON c.chat_id = cp.chat_id WHERE c.chat_id=$1`
	GetChatParticipantsQuery = `SELECT participant_id FROM chat_participants WHERE chat_id=$1`
	FetchChatListQuery       = `SELECT cp.chat_id, cp.chat_name, c.photo_url FROM chat_participants AS cp 
    LEFT JOIN chats AS c ON cp.chat_id = c.chat_id                      
                          WHERE cp.participant_id=$1`
	UpdatePhotoURLQuery     = "UPDATE chats SET photo_url=$1 WHERE chat_id=$2"
	RemoveUserFromChatQuery = "DELETE FROM chat_participants WHERE participant_id=$1 AND chat_id=$2"
	DeleteChatQuery         = "DELETE FROM chats WHERE chat_id=$1"
	DeleteMessageQuery      = "DELETE FROM messages WHERE msg_id=$1"
)

type PostgresRepo struct {
	pool *sql.DB
}

func NewPostgresRepo(pool *sql.DB) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

// GetChat
func (pr PostgresRepo) GetChat(ctx context.Context, chat models.Chat) (models.Chat, models.StatusCode) {
	row := pr.pool.QueryRowContext(ctx, GetChatInfoQuery, chat.ChatID)
	err := row.Scan(&chat.ChatID, &chat.Name, &chat.PhotoURL, &chat.CreatedAt)
	if err != nil {
		log.Println("err here get", err.Error())
		return models.Chat{}, models.NotFound
	}
	rows, err := pr.pool.QueryContext(ctx, GetChatParticipantsQuery, chat.ChatID)
	if err != nil {
		log.Println("err here get 2", err.Error())
		return models.Chat{}, models.InternalError
	}
	for rows.Next() {
		var participantID uuid.UUID
		err = rows.Scan(&participantID)
		if err != nil {
			slog.Error(err.Error())
		}
		chat.Participants = append(chat.Participants, participantID)
	}
	return chat, models.OK
}

// GetChatMessages
func (pr PostgresRepo) GetChatMessages(ctx context.Context, chat models.Chat, opts models.Opts) (models.Messages, models.StatusCode) {
	rows, err := pr.pool.QueryContext(ctx, GetChatMessagesQuery, chat.ChatID, opts.Page, opts.Limit)
	if err != nil {
		return nil, models.NotFound
	}

	msgs := make(models.Messages, 0)
	for rows.Next() {
		msg := models.Message{}
		err := rows.Scan(&msg.MsgID, &msg.SenderID, &msg.Payload, &msg.CreatedAt)
		if err != nil {
			return nil, models.InternalError
		}
		msg.ChatID = chat.ChatID
		msgs = append(msgs, msg)
	}
	sort.Sort(msgs)
	return msgs, models.OK
}

// FetchChatList
func (pr PostgresRepo) FetchChatList(ctx context.Context, user models.User) ([]models.ChatItem, models.StatusCode) {
	rows, err := pr.pool.QueryContext(ctx, FetchChatListQuery, user.ID)
	if err != nil {
		return nil, models.NotFound
	}

	chatList := make([]models.ChatItem, 0)
	for rows.Next() {
		chat := models.ChatItem{}
		err := rows.Scan(&chat.ChatID, &chat.Name, &chat.PhotoURL)
		if err != nil {
			return nil, models.InternalError
		}

		chatList = append(chatList, chat)
	}

	return chatList, models.OK
}

// CreateChat
func (pr PostgresRepo) CreateChat(ctx context.Context, chat models.Chat,
	chatNames map[string]string) models.StatusCode {
	tx, err := pr.pool.Begin()
	if err != nil {
		return models.InternalError
	}

	for _, participant := range chat.Participants {
		_, err := pr.pool.ExecContext(ctx, CreateChatParticipantsQuery, chat.ChatID, participant,
			chatNames[participant.String()])
		if err != nil {
			slog.Error("Failed to add a chat user-1", "user", participant.String())
			log.Println("err here pg", err.Error())
			txErr := tx.Rollback()
			if txErr != nil {
				slog.Error(txErr.Error())
			}
			return models.InternalError
		}
	}

	res, err := pr.pool.ExecContext(ctx, CreateChatQuery, chat.ChatID, chat.PhotoURL, chat.CreatedAt)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			slog.Error(txErr.Error())
		}
		return models.InternalError
	}
	if val, err := res.RowsAffected(); err != nil || val == 0 {
		txErr := tx.Rollback()
		if txErr != nil {
			slog.Error(txErr.Error())
		}
		return models.InternalError
	}

	txErr := tx.Commit()
	if txErr != nil {
		return models.InternalError
	}
	return models.OK
}

func (pr PostgresRepo) UpdateChatPhotoURL(ctx context.Context, chat models.Chat,
	photoURL string) models.StatusCode {
	res, err := pr.pool.ExecContext(ctx, UpdatePhotoURLQuery, photoURL, chat.ChatID)
	if err != nil {
		return models.InternalError
	}
	if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return models.InternalError
	}
	return models.OK
}

func (pr PostgresRepo) AddUsersToChat(ctx context.Context,
	chat models.Chat, chatNames map[string]string, users ...models.User) models.StatusCode {
	tx, err := pr.pool.Begin()
	if err != nil {
		return models.InternalError
	}
	if len(chat.Participants) > 0 {
		for _, user := range users {
			res, err := tx.ExecContext(ctx, CreateChatParticipantsQuery,
				chat.ChatID, user.ID, chatNames[user.ID.String()])
			if err != nil {
				slog.Error("Failed to add a chat user", "user", user.ID.String())
				txErr := tx.Rollback()
				if txErr != nil {
					slog.Error(txErr.Error())
				}
				return models.InternalError
			}
			if affected, err := res.RowsAffected(); err != nil || affected == 0 {
				txErr := tx.Rollback()
				if txErr != nil {
					slog.Error(txErr.Error())
				}
				return models.InternalError
			}
		}
	}
	txErr := tx.Commit()
	if txErr != nil {
		return models.InternalError
	}
	return models.OK
}

func (pr PostgresRepo) RemoveUserFromChat(ctx context.Context,
	chat models.Chat, users ...models.User) models.StatusCode {
	tx, err := pr.pool.Begin()
	if err != nil {
		return models.InternalError
	}
	if len(users) > 0 {
		for _, participant := range users {
			res, err := tx.ExecContext(ctx, RemoveUserFromChatQuery, participant.ID, chat.ChatID)
			if err != nil {
				slog.Error("Failed to remove a chat user", "user", participant.ID.String())
				txErr := tx.Rollback()
				if txErr != nil {
					slog.Error(txErr.Error())
				}
				return models.InternalError
			}
			if affected, err := res.RowsAffected(); err != nil || affected == 0 {
				txErr := tx.Rollback()
				if txErr != nil {
					slog.Error(txErr.Error())
				}
				return models.InternalError
			}
		}
	}
	txErr := tx.Commit()
	if txErr != nil {
		return models.InternalError
	}
	return models.OK
}

func (pr PostgresRepo) DeleteChat(ctx context.Context, chat models.Chat) models.StatusCode {
	_, err := pr.pool.ExecContext(ctx, DeleteChatQuery, chat.ChatID)
	if err != nil {
		return models.InternalError
	}
	return models.Deleted
}

func (pr PostgresRepo) DeleteMessage(ctx context.Context, message models.Message) models.StatusCode {
	_, err := pr.pool.ExecContext(ctx, DeleteMessageQuery, message.MsgID)
	if err != nil {
		return models.InternalError
	}
	return models.Deleted
}

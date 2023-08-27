package repo

import (
	"context"
	"github.com/google/uuid"
	"log"
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slog"
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
	pool internal.DB
}

func NewPostgresRepo(pool internal.DB) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

// GetChat
func (pr PostgresRepo) GetChat(chat models.Chat) (models.Chat, error) {
	ctx := context.Background()
	row := pr.pool.QueryRow(ctx, GetChatInfoQuery, chat.ChatID)
	err := row.Scan(&chat.ChatID, &chat.Name, &chat.PhotoURL, &chat.CreatedAt)
	if err != nil {
		log.Println("err here get", err.Error())
		return models.Chat{}, err
	}
	rows, err := pr.pool.Query(ctx, GetChatParticipantsQuery, chat.ChatID)
	if err != nil {
		log.Println("err here get 2", err.Error())
		return models.Chat{}, err
	}
	for rows.Next() {
		var participantID uuid.UUID
		err = rows.Scan(&participantID)
		if err != nil {
			slog.Error(err.Error())
		}
		chat.Participants = append(chat.Participants, participantID)
	}
	return chat, nil
}

// GetChatMessages
func (pr PostgresRepo) GetChatMessages(chat models.Chat, opts models.Opts) (models.Messages, error) {
	ctx := context.Background()
	rows, err := pr.pool.Query(ctx, GetChatMessagesQuery, chat.ChatID, opts.Page, opts.Limit)
	if err != nil {
		return nil, err
	}

	msgs := make(models.Messages, 0)
	for rows.Next() {
		msg := models.Message{}
		err := rows.Scan(&msg.MsgID, &msg.SenderID, &msg.Payload, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		msg.ChatID = chat.ChatID
		msgs = append(msgs, msg)
	}
	sort.Sort(msgs)
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
		err := rows.Scan(&chat.ChatID, &chat.Name, &chat.PhotoURL)
		if err != nil {
			return nil, err
		}

		chatList = append(chatList, chat)
	}

	log.Println("====== list from repo", chatList)

	return chatList, nil
}

// CreateChat
func (pr PostgresRepo) CreateChat(chat models.Chat,
	chatNames map[string]string) error {
	ctx := context.Background()
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, participant := range chat.Participants {
		batch.Queue(CreateChatParticipantsQuery, chat.ChatID, participant,
			chatNames[participant.String()])
	}

	batch.Queue(CreateChatQuery, chat.ChatID, chat.PhotoURL, chat.CreatedAt)

	results := pr.pool.SendBatch(ctx, batch)
	defer results.Close()
	for _, participant := range chat.Participants {
		_, err := results.Exec()
		if err != nil {
			slog.Error("Failed to add a chat user-1", "user", participant.String())
			log.Println("err here pg", err.Error())
			txErr := tx.Rollback(ctx)
			if err != nil {
				slog.Error(txErr.Error())
			}
			return err
		}
	}
	_, err = results.Exec()
	if err != nil {
		slog.Error("Failed to add a chat")
		txErr := tx.Rollback(ctx)
		if err != nil {
			slog.Error(txErr.Error())
		}
		return err
	}

	txErr := tx.Commit(ctx)
	if err != nil {
		slog.Error(txErr.Error())
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
				batch.Queue(CreateChatParticipantsQuery, chat.ChatID, participant)
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
				batch.Queue(RemoveUserFromChatQuery, chat.ChatID, participant)
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

func (pr PostgresRepo) DeleteChat(chat models.Chat) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := pr.pool.Exec(ctx, DeleteChatQuery, chat.ChatID)
	if err != nil {
		return err
	}
	return err
}

func (pr PostgresRepo) DeleteMessage(message models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := pr.pool.Exec(ctx, DeleteMessageQuery, message.MsgID)
	if err != nil {
		return err
	}
	return err
}

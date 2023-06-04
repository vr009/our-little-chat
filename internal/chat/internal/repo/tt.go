package repo

import (
	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"

	"github.com/tarantool/go-tarantool"
	"golang.org/x/exp/slog"
)

type TarantoolRepo struct {
	conn *tarantool.Connection
}

func NewTarantoolRepo(conn *tarantool.Connection) *TarantoolRepo {
	return &TarantoolRepo{conn: conn}
}

func (tt *TarantoolRepo) GetFreshMessagesFromChat(chat models2.Chat) ([]models.Message, error) {
	conn := tt.conn
	var msgs []models.Message
	err := conn.CallTyped("fetch_msgs", []interface{}{chat.ChatID}, &msgs)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (tt *TarantoolRepo) InsertChat(chat models2.Chat) (models2.Chat, error) {
	conn := tt.conn
	participants := []string{}
	for _, paricipant := range chat.Participants {
		participants = append(participants, paricipant.String())
	}
	slog.Info("to insert", "list", participants, "list uuids", chat.Participants)
	_, err := conn.Call("create_chat", []interface{}{participants, chat.ChatID.String()})
	if err != nil {
		return models2.Chat{}, err
	}
	return chat, nil
}

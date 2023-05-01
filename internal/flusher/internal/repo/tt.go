package repo

import (
	"our-little-chatik/internal/models"

	"github.com/tarantool/go-tarantool"
)

type TarantoolRepo struct {
	conn *tarantool.Connection
}

func NewTarantoolRepo(conn *tarantool.Connection) *TarantoolRepo {
	return &TarantoolRepo{conn: conn}
}

func (tt *TarantoolRepo) FetchAllMessages() ([]models.Message, error) {
	conn := tt.conn
	msgs := []models.Message{}
	err := conn.CallTyped("flush_msgs", []interface{}{}, &msgs)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (tt *TarantoolRepo) FetchChatListUpdate() ([]models.ChatItem, error) {
	conn := tt.conn
	chats := []models.ChatItem{}
	err := conn.CallTyped("fetch_all_chats_upd", []interface{}{}, &chats)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (tt *TarantoolRepo) FetchChatParticipants() ([]models.Chat, error) {
	conn := tt.conn
	chats := []models.Chat{}
	err := conn.CallTyped("flush_chats_participants", []interface{}{}, &chats)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

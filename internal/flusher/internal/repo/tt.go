package repo

import (
	"github.com/tarantool/go-tarantool"
	"our-little-chatik/internal/models"
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

func (tt *TarantoolRepo) FetchAllChats() ([]models.Chat, error) {
	conn := tt.conn
	chats := []models.Chat{}
	err := conn.CallTyped("flush_chats", []interface{}{}, &chats)
	if err != nil {
		return nil, err
	}
	return chats, nil
}
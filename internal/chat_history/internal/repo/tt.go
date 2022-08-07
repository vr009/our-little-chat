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

func (tt *TarantoolRepo) GetFreshChat(chat models.Chat) ([]models.Message, error) {
	conn := tt.conn
	msgs := []models.Message{}
	err := conn.CallTyped("fetch", []interface{}{chat.ChatID}, &msgs)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

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

func (tt *TarantoolRepo) FetchAll() ([]models.Message, error) {
	conn := tt.conn
	msgs := []models.Message{}
	err := conn.CallTyped("flush", []interface{}{}, &msgs)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

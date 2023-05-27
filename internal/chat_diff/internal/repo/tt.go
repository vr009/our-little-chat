package repo

import (
	"fmt"

	"our-little-chatik/internal/models"

	"github.com/tarantool/go-tarantool"
)

// TODO rewrite it to websockets
// See https://github.com/tarantool/websocket

type TarantoolRepo struct {
	conn *tarantool.Connection
}

func NewTarantoolRepo(conn *tarantool.Connection) *TarantoolRepo {
	return &TarantoolRepo{conn: conn}
}

func (tt *TarantoolRepo) FetchUpdates(user models.User) []models.ChatItem {
	conn := tt.conn
	updates := []models.ChatItem{}

	err := conn.CallTyped("fetch_unread_messages", []interface{}{user.UserID.String()}, &updates)
	if err != nil {
		return nil
	}

	fmt.Println("updates", updates)
	return updates
}

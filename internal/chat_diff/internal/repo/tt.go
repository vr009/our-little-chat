package repo

import (
	"fmt"

	"our-little-chatik/internal/models"

	"github.com/google/uuid"
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

	resp, err := conn.Call("fetch_unread_messages", []interface{}{user.UserID.String()})
	if err != nil {
		return nil
	}

	updates = make([]models.ChatItem, len(resp.Data))
	for i, el := range resp.Data {
		sl := el.([]interface{})
		if len(sl) < 1 {
			return nil
		}
		updates[i].ChatID, _ = uuid.Parse(sl[0].(string))
		updates[i].LastSender, _ = uuid.Parse(sl[1].(string))
		updates[i].LastMsg = sl[2].(string)
	}

	fmt.Println(updates)
	return updates
}

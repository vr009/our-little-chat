package repo

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"
	models2 "our-little-chatik/internal/chat_diff/internal/models"
)

// TODO rewrite it to websockets
// See https://github.com/tarantool/websocket

type TarantoolRepo struct {
	conn *tarantool.Connection
}

func NewTarantoolRepo(conn *tarantool.Connection) *TarantoolRepo {
	return &TarantoolRepo{conn: conn}
}

func (tt *TarantoolRepo) FetchUpdates(user models2.ChatUser) []models2.ChatItem {
	conn := tt.conn
	updates := []models2.ChatItem{}

	//err := conn.CallTyped("fetch_chats_upd", []interface{}{user.ID}, &updates)
	resp, err := conn.Call("fetch_chats_upd", []interface{}{user.ID})
	if err != nil {
		return nil
	}

	updates = make([]models2.ChatItem, len(resp.Data))
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

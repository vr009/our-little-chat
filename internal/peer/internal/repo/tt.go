package repo

import (
	"errors"
	"fmt"
	"log"

	"our-little-chatik/internal/peer/internal/models"

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

func (tt *TarantoolRepo) SendPayload(msg *models.Message) error {
	conn := tt.conn
	resp, err := conn.Call("put", []interface{}{
		msg.ChatID.String(),
		msg.SenderID.String(),
		msg.Payload})
	if err != nil {
		fmt.Println(err)
		return err
	}
	if resp == nil {
		fmt.Println("Response is nil after Call")
		return errors.New("Response is nil after Call")
	}
	if len(resp.Data) < 1 {
		return errors.New("Response.Data is empty after Eval")
	}
	return nil
}
func (tt *TarantoolRepo) FetchUpdates(chat *models.Chat, peer *models.Peer) ([]models.Message, error) {
	var msgs []models.Message
	conn := tt.conn
	err := conn.CallTyped("take_msgs",
		[]interface{}{chat.ChatID.String(),
			peer.PeerID.String()}, &msgs)
	if err != nil && len(msgs) < 1 {
		log.Println("error from tt: ", err)
		return nil, err
	}

	if len(msgs) > 0 && msgs[0].Payload == "" {
		return nil, nil
	}
	fmt.Printf("fetched for %s\n", peer.PeerID.String())

	return msgs, nil
}

func (tt *TarantoolRepo) Close() {
	tt.conn.Close()
}

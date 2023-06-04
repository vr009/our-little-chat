package repo

import (
	"errors"

	"our-little-chatik/internal/peer/internal/models"

	"github.com/tarantool/go-tarantool"
	"golang.org/x/exp/slog"
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
		slog.Error(err.Error())
		return err
	}
	if resp == nil {
		slog.Error("Response is nil after Call")
		return errors.New("Response is nil after Call")
	}
	if len(resp.Data) < 1 {
		slog.Error("Response.Data is empty")
		return errors.New("Response.Data is empty")
	}
	slog.Info("SUCCEED")
	return nil
}
func (tt *TarantoolRepo) FetchUpdates(chat *models.Chat, peer *models.Peer) ([]models.Message, error) {
	var msgs []models.Message
	conn := tt.conn
	err := conn.CallTyped("take_msgs",
		[]interface{}{chat.ChatID.String(),
			peer.PeerID.String()}, &msgs)
	if err != nil && len(msgs) < 1 {
		slog.Error("error from tt: ", err)
		return nil, err
	}

	if len(msgs) > 0 && msgs[0].Payload == "" {
		return nil, nil
	}
	slog.Info("fetched for " + peer.PeerID.String())

	return msgs, nil
}

func (tt *TarantoolRepo) Close() {
	tt.conn.Close()
}

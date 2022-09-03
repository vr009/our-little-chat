package test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tarantool/go-tarantool"
	"gopkg.in/vmihailenco/msgpack.v2"
	"testing"
	"time"
)

var server = "127.0.0.1:3301"
var opts = tarantool.Opts{
	Timeout: 500 * time.Millisecond,
	User:    "test",
	Pass:    "test",
	//Concurrency: 32,
	//RateLimit: 4*1024,
}

func TestPeers(t *testing.T) {
	var err error

	conn, err := tarantool.Connect(server, opts)
	if err != nil {
		t.Error("No connection available:", err)
		return
	}
	defer conn.Close()

	userA := uuid.New()
	userB := uuid.New()
	chatID := uuid.New()

	msg1 := "hi"
	msg2 := "hey"

	resp, err := conn.Call("put", []interface{}{chatID.String(), userA.String(), userB.String(), msg1})
	if err != nil {
		t.Fatalf("Failed to Call: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Response is nil after Call")
	}
	if len(resp.Data) < 1 {
		t.Errorf("Response.Data is empty after Eval")
	}
	t.Log(resp.Data)

	resp, err = conn.Call("put", []interface{}{chatID.String(), userB.String(), userA.String(), msg2})
	if err != nil {
		t.Fatalf("Failed to Call: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Response is nil after Call")
	}
	if len(resp.Data) < 1 {
		t.Errorf("Response.Data is empty after Eval")
	}
	t.Log(resp.Data)

	msgs := []Message{}

	err = conn.CallTyped("take_msgs", []interface{}{chatID.String(), userB.String()}, &msgs)
	if err != nil {
		t.Fatalf("Failed to Call: %s", err.Error())
	}
	t.Log(resp.Data)

	resp, err = conn.Call("fetch_chats_upd", []interface{}{userB.String()})
	if err != nil {
		t.Fatalf("Failed to Call: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Response is nil after Call")
	}
	if len(resp.Data) < 1 {
		t.Errorf("Response.Data is empty after Eval")
	}
	t.Log(resp.Data)

}

type Message struct {
	ChatID    uuid.UUID `json:"chatID"`
	SenderID  uuid.UUID `json:"senderID"`
	MsgID     uuid.UUID `json:"-"`
	Payload   string    `json:"payload"`
	CreatedAt float64   `json:"-"`
}

func (m *Message) EncodeMsgpack(e *msgpack.Encoder) error {
	return nil
}

func (m *Message) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var uuidStr string
	var l int
	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}
	if l != 5 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	//chat id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	m.ChatID, _ = uuid.Parse(uuidStr)
	//msg id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	m.MsgID, _ = uuid.Parse(uuidStr)
	//sender id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	m.SenderID, _ = uuid.Parse(uuidStr)
	//payload
	if m.Payload, err = d.DecodeString(); err != nil {
		return err
	}
	//timestamp
	if m.CreatedAt, err = d.DecodeFloat64(); err != nil {
		return err
	}
	return nil
}

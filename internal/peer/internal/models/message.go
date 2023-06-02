package models

import (
	"fmt"

	"github.com/google/uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Message struct {
	ChatID       uuid.UUID `json:"chat_id"`
	SenderID     uuid.UUID `json:"sender_id"`
	MsgID        uuid.UUID `json:"-"`
	Payload      string    `json:"payload"`
	CreatedAt    float64   `json:"-"`
	SessionStart bool      `json:"session_start,omitempty"`
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

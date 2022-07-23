package models

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type ChatItem struct {
	ChatID      uuid.UUID `json:"chat_id" bson:"chat_id"`
	LastSender  uuid.UUID `json:"last_sender" bson:"last_sender"`
	LastMsg     string    `json:"last_msg" bson:"last_msg"`
	LastMessage float64   `json:"last_message" bson:"last_message"`
}

func (ch *ChatItem) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var uuidStr string
	var l int
	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}
	if l != 4 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}

	//chat id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	ch.ChatID, _ = uuid.Parse(uuidStr)
	//sender id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	ch.LastSender, _ = uuid.Parse(uuidStr)
	//payload
	if ch.LastMsg, err = d.DecodeString(); err != nil {
		return err
	}

	ch.LastMessage, _ = d.DecodeFloat64()
	//timestamp
	return nil
}
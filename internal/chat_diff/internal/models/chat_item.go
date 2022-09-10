package models

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
	"time"
)

type ChatItem struct {
	ChatID     uuid.UUID `json:"chatId"`
	LastSender uuid.UUID `json:"lastSender"`
	LastMsg    string    `json:"lastMsg"`
	LastUpdate time.Time `json:"lastUpdate"`
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

	ch.LastUpdate, _ = d.DecodeTime()
	//timestamp
	return nil
}

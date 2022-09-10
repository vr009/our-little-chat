package models

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Chat struct {
	Participant uuid.UUID `json:"participant" bson:"participant"`
	ChatID      uuid.UUID `json:"chat_id" bson:"chat_id"`
	LastRead    float64   `json:"last_read" bson:"last_read"`
}

func (c *Chat) EncodeMsgpack(e *msgpack.Encoder) error {
	return nil
}

func (c *Chat) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var uuidStr string
	var l int
	if l, err = d.DecodeSliceLen(); err != nil {
		return err
	}
	if l != 3 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	//participant id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}

	id, _ := uuid.Parse(uuidStr)
	c.Participant = id

	//chat id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	c.ChatID, _ = uuid.Parse(uuidStr)

	//timestamp
	if c.LastRead, err = d.DecodeFloat64(); err != nil {
		return err
	}
	return nil
}

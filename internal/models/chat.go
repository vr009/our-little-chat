package models

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Chat struct {
	ChatID       uuid.UUID   `json:"chat_id" bson:"chat_id"`
	Participants []uuid.UUID `json:"participants" bson:"participants"`
	LastMessage  float64     `json:"last_message" bson:"last_message"`
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
	if l != 4 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	//chat id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	c.ChatID, _ = uuid.Parse(uuidStr)

	//sender id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}

	id, _ := uuid.Parse(uuidStr)
	c.Participants = append(c.Participants, id)

	//receiver id
	if uuidStr, err = d.DecodeString(); err != nil {
		return err
	}
	id, _ = uuid.Parse(uuidStr)
	c.Participants = append(c.Participants, id)

	//timestamp
	//if c.LastUpdate, err = d.DecodeTime(); err != nil {
	//	return err
	//}
	if c.LastMessage, err = d.DecodeFloat64(); err != nil {
		return err
	}
	return nil
}

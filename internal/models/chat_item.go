package models

import (
	"github.com/google/uuid"
)

type ChatItem struct {
	ChatID          uuid.UUID `json:"chat_id"`
	LastSender      uuid.UUID `json:"last_sender"`
	MsgID           uuid.UUID `json:"msg_id"`
	LastMsg         string    `json:"last_msg"`
	LastMessageTime int64     `json:"last_message"`
}

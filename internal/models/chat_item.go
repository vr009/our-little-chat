package models

import (
	"database/sql"
	"github.com/google/uuid"
)

type ChatItem struct {
	ChatID          uuid.UUID      `json:"chat_id"`
	Name            string         `json:"name,omitempty"`
	PhotoURL        string         `json:"photo_url,omitempty"`
	LastMsg         sql.NullString `json:"last_msg"`
	LastMessageTime sql.NullInt64  `json:"last_message"`
}

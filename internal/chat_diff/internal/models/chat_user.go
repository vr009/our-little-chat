package models

import (
	"github.com/google/uuid"
)

type ChatUser struct {
	ID       uuid.UUID
	Username string
	Updates  chan []ChatItem
}

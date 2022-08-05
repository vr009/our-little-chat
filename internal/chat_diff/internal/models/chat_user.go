package models

import (
	"github.com/google/uuid"
)

type ChatUser struct {
	ID       uuid.UUID
	Username string
	Updates  chan []ChatItem
}

//
//func NewChatUserFromAuth(auth *Auth, resolver internal.TokenResolver) *ChatUser {
//	id, err := resolver.ResolveToken(auth.Token)
//	if err != nil {
//		return nil
//	}
//	return &ChatUser{ID: id}
//}

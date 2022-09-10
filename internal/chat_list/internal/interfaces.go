package internal

import "our-little-chatik/internal/models"

type Usecase interface {
	GetChatList(user models.User) ([]models.ChatItem, error)
}

type Repo interface {
	FetchChatList(user models.User) ([]models.ChatItem, error)
}

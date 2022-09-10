package usecase

import (
	"our-little-chatik/internal/chat_list/internal"
	"our-little-chatik/internal/models"
)

type ChatListUsecase struct {
	repo internal.Repo
}

func NewChatListUsecase(repo internal.Repo) *ChatListUsecase {
	return &ChatListUsecase{
		repo: repo,
	}
}

func (clu *ChatListUsecase) GetChatList(user models.User) ([]models.ChatItem, error) {
	return clu.repo.FetchChatList(user)
}

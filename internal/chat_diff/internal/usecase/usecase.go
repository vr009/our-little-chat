package usecase

import (
	"our-little-chatik/internal/chat_diff/internal"
	models2 "our-little-chatik/internal/chat_diff/internal/models"
)

type Usecase struct {
	repo internal.ChatDiffRepo
}

func NewUsecase(repo internal.ChatDiffRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) GetUpdates(user models2.ChatUser) []models2.ChatItem {
	return uc.repo.FetchUpdates(user)
}

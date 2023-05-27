package usecase

import (
	"our-little-chatik/internal/chat_diff/internal"
	"our-little-chatik/internal/models"
)

type Usecase struct {
	repo internal.ChatDiffRepo
}

func NewUsecase(repo internal.ChatDiffRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) GetUpdates(user models.User) []models.ChatItem {
	return uc.repo.FetchUpdates(user)
}

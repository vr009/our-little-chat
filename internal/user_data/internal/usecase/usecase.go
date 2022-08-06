package usecase

import (
	"user_data/internal"
	"user_data/internal/models"
)

type UserdataUseCase struct {
	repo internal.UserdataRepo
}

func NewUserdataUseCase(base internal.UserdataRepo) *UserdataUseCase {
	return &UserdataUseCase{
		repo: base,
	}
}

func (uc *UserdataUseCase) GetAllUsers(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.GetAllUsers(userData)
}

func (uc *UserdataUseCase) CreateUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.CreateUser(userData)
}

func (uc *UserdataUseCase) GetUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.GetUser(userData)
}

func (uc *UserdataUseCase) DeleteUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.GetUser(userData)
}

func (uc *UserdataUseCase) UpdateUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.GetUser(userData)
}

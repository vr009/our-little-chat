package usecase

import (
	"github.com/google/uuid"
	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"
	"time"
)

type UserdataUseCase struct {
	repo internal.UserdataRepo
}

func NewUserdataUseCase(base internal.UserdataRepo) *UserdataUseCase {
	return &UserdataUseCase{
		repo: base,
	}
}

func (uc *UserdataUseCase) GetAllUsers() ([]models.UserData, models.StatusCode) {
	return uc.repo.GetAllUsers()
}

func (uc *UserdataUseCase) CreateUser(userData models.UserData) (models.UserData, models.StatusCode) {
	userData.UserID = uuid.New()
	userData.Registered = time.Now()
	return uc.repo.CreateUser(userData)
}

func (uc *UserdataUseCase) GetUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.GetUser(userData)
}

func (uc *UserdataUseCase) DeleteUser(userData models.UserData) models.StatusCode {
	return uc.repo.DeleteUser(userData)
}

func (uc *UserdataUseCase) UpdateUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.UpdateUser(userData)
}

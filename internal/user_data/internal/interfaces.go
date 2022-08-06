package internal

import (
	"user_data/internal/models"
)

type UserdataRepo interface {
	GetAllUsers(userData models.UserData) (models.UserData, models.StatusCode)
	CreateUser(userData models.UserData) (models.UserData, models.StatusCode)
	GetUser(userData models.UserData) (models.UserData, models.StatusCode)
	DeleteUser(userData models.UserData) (models.UserData, models.StatusCode)
	UpdateUser(userData models.UserData) (models.UserData, models.StatusCode)
}

type UserdataUseCase interface {
	GetAllUsers(userData models.UserData) (models.UserData, models.StatusCode)
	CreateUser(userData models.UserData) models.UserData
	GetUser(userData models.UserData) (models.UserData, models.StatusCode)
	DeleteUser(userData models.UserData) (models.UserData, models.StatusCode)
	UpdateUser(userData models.UserData) (models.UserData, models.StatusCode)
}

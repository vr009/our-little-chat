package internal

import "our-little-chatik/internal/user_data/internal/models"

type UserdataRepo interface {
	GetAllUsers() ([]models.UserData, models.StatusCode)
	CreateUser(userData models.UserData) (models.UserData, models.StatusCode)
	GetUser(userData models.UserData) (models.UserData, models.StatusCode)
	GetUserForItsName(userData models.UserData) (models.UserData, models.StatusCode)
	DeleteUser(userData models.UserData) models.StatusCode
	UpdateUser(personNew models.UserData) (models.UserData, models.StatusCode)
	FindUser(name string) ([]models.UserData, models.StatusCode)
}

type UserdataUseCase interface {
	GetAllUsers() ([]models.UserData, models.StatusCode)
	CreateUser(userData models.UserData) (models.UserData, models.StatusCode)
	GetUser(userData models.UserData) (models.UserData, models.StatusCode)
	DeleteUser(userData models.UserData) models.StatusCode
	UpdateUser(userData models.UserData) (models.UserData, models.StatusCode)
	CheckUser(userData models.UserData) (models.UserData, models.StatusCode)
	FindUser(name string) ([]models.UserData, models.StatusCode)
}

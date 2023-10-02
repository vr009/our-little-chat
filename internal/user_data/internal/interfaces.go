package internal

import (
	models2 "our-little-chatik/internal/models"
)

type UserRepo interface {
	GetAllUsers() ([]models2.User, models2.StatusCode)
	CreateUser(User models2.User) (models2.User, models2.StatusCode)
	GetUser(User models2.User) (models2.User, models2.StatusCode)
	GetUserForItsName(User models2.User) (models2.User, models2.StatusCode)
	DeleteUser(User models2.User) models2.StatusCode
	UpdateUser(personNew models2.User) (models2.User, models2.StatusCode)
	FindUser(name string) ([]models2.User, models2.StatusCode)
}

type UserUsecase interface {
	GetAllUsers() ([]models2.User, models2.StatusCode)
	CreateUser(User models2.User) (models2.User, models2.StatusCode)
	GetUser(User models2.User) (models2.User, models2.StatusCode)
	DeleteUser(User models2.User) models2.StatusCode
	UpdateUser(User models2.User) (models2.User, models2.StatusCode)
	CheckUser(User models2.User) (models2.User, models2.StatusCode)
	FindUser(name string) ([]models2.User, models2.StatusCode)
}

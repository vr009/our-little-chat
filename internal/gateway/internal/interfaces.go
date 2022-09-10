package internal

import (
	"gateway/internal/models"
)

type ServerUsecase interface {
	SignUp(user *models.User) (*models.User, models.ErrorCode)
	SignIn(user *models.User) (*models.User, models.ErrorCode)

	GetSessionID(user models.SessionUser) (models.SessionUser, error)
	AddSession(user models.SessionUser) (models.SessionUser, error)
	RemoveSession(user models.SessionUser) error
}

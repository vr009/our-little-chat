package internal

import (
	"our-little-chatik/internal/models"
)

type GatewayUsecase interface {
	SignUp(user *models.User) (*models.Session, error)
	SignIn(user *models.User) (*models.Session, error)

	GetSessionFromUser(user models.User) (*models.Session, error)
	GetUserFromSession(user models.Session) (*models.User, error)
}

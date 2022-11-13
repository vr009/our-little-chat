package internal

import (
	models2 "our-little-chatik/internal/auth/internal/models"
	"our-little-chatik/internal/models"
)

type AuthRepo interface {
	CreateSession(session models.Session) (models.Session, models2.StatusCode)
	DeleteSession(session models.Session) models2.StatusCode
	GetToken(session models.Session) (models.Session, models2.StatusCode)
	GetUser(session models.Session) (models.Session, models2.StatusCode)
}

type AuthUseCase interface {
	CreateSession(session models.Session) (models.Session, models2.StatusCode)
	DeleteSession(session models.Session) models2.StatusCode
	GetToken(session models.Session) (models.Session, models2.StatusCode)
	GetUser(session models.Session) (models.Session, models2.StatusCode)
}

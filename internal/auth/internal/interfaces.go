package internal

import "our-little-chatik/internal/auth/internal/models"

type AuthRepo interface {
	CreateSession(session models.Session) (models.Session, models.StatusCode)
	DeleteSession(session models.Session) models.StatusCode
	GetToken(session models.Session) (models.Session, models.StatusCode)
	GetUser(session models.Session) (models.Session, models.StatusCode)
}

type AuthUseCase interface {
	CreateSession(session models.Session) (models.Session, models.StatusCode)
	DeleteSession(session models.Session) models.StatusCode
	GetToken(session models.Session) (models.Session, models.StatusCode)
	GetUser(session models.Session) (models.Session, models.StatusCode)
}

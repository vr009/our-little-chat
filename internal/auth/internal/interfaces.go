package internal

import "auth/internal/models"

type AuthRepo interface {
	CreateSession(session models.Session) models.Session
	DeleteSession(session models.Session)
	GetSession(session models.Session) models.Session
}

type AuthUseCase interface {
	CreateSession(session models.Session) models.Session
	DeleteSession(session models.Session)
	GetSession(session models.Session) models.Session
}

package usecase

import (
	"auth/internal"
	"auth/internal/models"
)

type AuthUseCase struct {
	repo internal.AuthRepo
}

func NewAuthUseCase(base internal.AuthRepo) *AuthUseCase {
	return &AuthUseCase{
		repo: base,
	}
}

func (uc *AuthUseCase) CreateSession(session models.Session) models.Session {
	return uc.repo.CreateSession(session)
}

func (uc *AuthUseCase) DeleteSession(session models.Session) {
	uc.repo.CreateSession(session)
}

func (uc *AuthUseCase) GetSession(session models.Session) models.Session {
	return uc.repo.GetSession(session)
}

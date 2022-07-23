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

func (uc *AuthUseCase) CreateSession(session models.Session) (models.Session, models.StatusCode) {
	return uc.repo.CreateSession(session)
}

func (uc *AuthUseCase) DeleteSession(session models.Session) models.StatusCode {
	return uc.repo.DeleteSession(session)
}

func (uc *AuthUseCase) GetToken(session models.Session) (models.Session, models.StatusCode) {
	return uc.repo.GetToken(session)
}

func (uc *AuthUseCase) GetUser(session models.Session) (models.Session, models.StatusCode) {
	return uc.repo.GetUser(session)
}

package usecase

import (
	"auth/internal"
	"auth/internal/models"
	"crypto/md5"
	"fmt"
)

type AuthUseCase struct {
	repo internal.AuthRepo
}

func NewAuthUseCase(base internal.AuthRepo) *AuthUseCase {
	return &AuthUseCase{
		repo: base,
	}
}

func MD5(data string) string {
	h := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", h)
}

func (uc *AuthUseCase) CreateSession(session models.Session) (models.Session, models.StatusCode) {
	session.Token = MD5(session.UserID.String())
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

package usecase

import (
	"crypto/md5"
	"fmt"

	"our-little-chatik/internal/auth/internal"
	models2 "our-little-chatik/internal/auth/internal/models"
	"our-little-chatik/internal/models"
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

func (uc *AuthUseCase) CreateSession(session models.Session) (models.Session, models2.StatusCode) {
	session.Token = MD5(session.UserID.String())
	return uc.repo.CreateSession(session)
}

func (uc *AuthUseCase) DeleteSession(session models.Session) models2.StatusCode {
	return uc.repo.DeleteSession(session)
}

func (uc *AuthUseCase) GetToken(session models.Session) (models.Session, models2.StatusCode) {
	return uc.repo.GetToken(session)
}

func (uc *AuthUseCase) GetUser(session models.Session) (models.Session, models2.StatusCode) {
	return uc.repo.GetUser(session)
}

package usecase

import (
	"golang.org/x/exp/slog"
	"time"

	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserdataUseCase struct {
	repo internal.UserdataRepo
}

func NewUserdataUseCase(base internal.UserdataRepo) *UserdataUseCase {
	return &UserdataUseCase{
		repo: base,
	}
}

func (uc *UserdataUseCase) GetAllUsers() ([]models.UserData, models.StatusCode) {
	return uc.repo.GetAllUsers()
}

func (uc *UserdataUseCase) CreateUser(userData models.UserData) (models.UserData, models.StatusCode) {
	pswd, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 10)
	if err != nil {
		return models.UserData{}, models.InternalError
	}
	userData.Password = string(pswd)
	userData.UserID = uuid.New()
	userData.Registered = time.Now()
	slog.Info("CREATE", "pswd", userData.Password, "id", userData.UserID)
	return uc.repo.CreateUser(userData)
}

func (uc *UserdataUseCase) GetUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.GetUser(userData)
}

func (uc *UserdataUseCase) DeleteUser(userData models.UserData) models.StatusCode {
	return uc.repo.DeleteUser(userData)
}

func (uc *UserdataUseCase) UpdateUser(personNew models.UserData) (models.UserData, models.StatusCode) {
	personOld, status := uc.GetUser(personNew)
	if status != models.OK {
		return models.UserData{}, status
	}

	if personNew.Name == "" {
		personNew.Name = personOld.Name
	}
	if personNew.Surname == "" {
		personNew.Surname = personOld.Surname
	}
	if personNew.Nickname == "" {
		personNew.Nickname = personOld.Nickname
	}
	return uc.repo.UpdateUser(personNew)
}

func (uc *UserdataUseCase) CheckUser(userData models.UserData) (models.UserData, models.StatusCode) {
	userFromRepo, code := uc.repo.GetUserForItsName(userData)
	if code == models.NotFound {
		return models.UserData{}, code
	}
	slog.Info("CHECK", "pswd", userData.Password, "id", userFromRepo.UserID)
	slog.Info("CHECK", "pswd", userFromRepo.Password, "id", userFromRepo.UserID)
	err := bcrypt.CompareHashAndPassword([]byte(userFromRepo.Password), []byte(userData.Password))
	if err != nil {
		slog.Error(err.Error())
		return models.UserData{}, models.Forbidden
	}
	return userFromRepo, models.OK
}

func (uc *UserdataUseCase) FindUser(name string) ([]models.UserData, models.StatusCode) {
	return uc.repo.FindUser(name)
}

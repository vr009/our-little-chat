package usecase

import (
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"our-little-chatik/internal/user_data/internal"
)

type UserdataUsecase struct {
	repo internal.UserdataRepo
}

func NewUserdataUsecase(base internal.UserdataRepo) *UserdataUsecase {
	return &UserdataUsecase{
		repo: base,
	}
}

func (uc *UserdataUsecase) GetAllUsers() ([]models.UserData, models.StatusCode) {
	return uc.repo.GetAllUsers()
}

func (uc *UserdataUsecase) CreateUser(userData models.UserData) (models.UserData, models.StatusCode) {
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

func (uc *UserdataUsecase) GetUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return uc.repo.GetUser(userData)
}

func (uc *UserdataUsecase) DeleteUser(userData models.UserData) models.StatusCode {
	return uc.repo.DeleteUser(userData)
}

func (uc *UserdataUsecase) UpdateUser(personNew models.UserData) (models.UserData, models.StatusCode) {
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

func (uc *UserdataUsecase) CheckUser(userData models.UserData) (models.UserData, models.StatusCode) {
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

func (uc *UserdataUsecase) FindUser(name string) ([]models.UserData, models.StatusCode) {
	return uc.repo.FindUser(name)
}

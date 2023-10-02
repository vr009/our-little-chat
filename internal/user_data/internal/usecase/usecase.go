package usecase

import (
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/models"
	"time"

	"github.com/google/uuid"
	"our-little-chatik/internal/user_data/internal"
)

type UserUsecase struct {
	repo internal.UserRepo
}

func NewUserUsecase(base internal.UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: base,
	}
}

func (uc *UserUsecase) GetAllUsers() ([]models.User, models.StatusCode) {
	return uc.repo.GetAllUsers()
}

func (uc *UserUsecase) CreateUser(User models.User) (models.User, models.StatusCode) {
	User.ID = uuid.New()
	User.Registered = time.Now()
	slog.Info("CREATE", "pswd", User.Password, "id", User.ID)
	return uc.repo.CreateUser(User)
}

func (uc *UserUsecase) GetUser(User models.User) (models.User, models.StatusCode) {
	return uc.repo.GetUser(User)
}

func (uc *UserUsecase) DeleteUser(User models.User) models.StatusCode {
	return uc.repo.DeleteUser(User)
}

func (uc *UserUsecase) UpdateUser(personNew models.User) (models.User, models.StatusCode) {
	personOld, status := uc.GetUser(personNew)
	if status != models.OK {
		return models.User{}, status
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
	if personNew.Password.Hash == nil {
		personNew.Nickname = personOld.Nickname
	}
	return uc.repo.UpdateUser(personNew)
}

func (uc *UserUsecase) CheckUser(User models.User) (models.User, models.StatusCode) {
	userFromRepo, code := uc.repo.GetUserForItsName(User)
	if code == models.NotFound {
		return models.User{}, code
	}
	return userFromRepo, models.OK
}

func (uc *UserUsecase) FindUser(name string) ([]models.User, models.StatusCode) {
	return uc.repo.FindUser(name)
}

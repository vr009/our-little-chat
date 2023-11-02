package usecase

import (
	"our-little-chatik/internal/models"
	models2 "our-little-chatik/internal/users/internal/models"
	"time"

	"github.com/google/uuid"
	"our-little-chatik/internal/users/internal"
)

type UserUsecase struct {
	repo internal.UserRepo
}

func NewUserUsecase(repo internal.UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uc *UserUsecase) SignUp(request models2.SignUpPersonRequest) (models.User, models.StatusCode) {
	user := models.User{
		Name:      *request.Name,
		Nickname:  *request.Nickname,
		Surname:   *request.Surname,
		Avatar:    *request.Avatar,
		Activated: true,
	}
	err := user.Password.Set(*request.Password)
	if err != nil {
		return models.User{}, models.InternalError
	}
	user.ID = uuid.New()
	user.RegisteredAt = time.Now()
	return uc.repo.CreateUser(user)
}

func (uc *UserUsecase) Login(request models2.LoginRequest) (models.User, models.StatusCode) {
	user, status := uc.repo.GetUserForItsNickname(models.User{Nickname: *request.Nickname})
	if status != models.OK {
		return models.User{}, status
	}

	if !user.Activated {
		return models.User{}, models.InActivated
	}

	ok, err := user.Password.Matches(*request.Password)
	if err != nil {
		return models.User{}, models.InternalError
	}
	if !ok {
		return models.User{}, models.Unauthorized
	}

	return user, models.OK
}

func (uc *UserUsecase) GetUser(request models2.GetUserRequest) (models.User, models.StatusCode) {
	return uc.repo.GetUserForItsID(models.User{ID: request.UserID})
}

func (uc *UserUsecase) DeactivateUser(user models.User) models.StatusCode {
	return uc.repo.DeactivateUser(user)
}

func (uc *UserUsecase) UpdateUser(userToUpdate models.User,
	request models2.UpdateUserRequest) (models.User, models.StatusCode) {
	oldUser, status := uc.repo.GetUserForItsID(userToUpdate)
	if status != models.OK {
		return models.User{}, status
	}

	if !oldUser.Activated {
		return models.User{}, models.InActivated
	}

	newUser := models.User{
		ID:        oldUser.ID,
		Activated: oldUser.Activated,
		Nickname:  oldUser.Nickname,
		Name:      oldUser.Name,
		Surname:   oldUser.Surname,
		Avatar:    oldUser.Avatar,
	}
	if request.Name != nil {
		newUser.Name = *request.Name
	}
	if request.Surname != nil {
		newUser.Surname = *request.Surname
	}
	if request.Nickname != nil {
		newUser.Nickname = *request.Nickname
	}
	if request.Avatar != nil {
		newUser.Avatar = *request.Avatar
	}
	if request.NewPassword != nil {
		match, err := oldUser.Password.Matches(*request.OldPassword)
		if err != nil {
			return models.User{}, models.InternalError
		}
		if !match {
			return models.User{}, models.Forbidden
		}
		err = newUser.Password.Set(*request.NewPassword)
		if err != nil {
			return models.User{}, models.InternalError
		}
	} else {
		newUser.Password.Hash = oldUser.Password.Hash
	}
	return uc.repo.UpdateUser(newUser)
}

func (uc *UserUsecase) FindUsers(name string) ([]models.User, models.StatusCode) {
	return uc.repo.FindUsers(name)
}

package internal

import (
	internalmodels "our-little-chatik/internal/models"
	"our-little-chatik/internal/users/internal/models"
)

type UserRepo interface {
	CreateUser(user internalmodels.User) (internalmodels.User, internalmodels.StatusCode)
	GetUserForItsID(user internalmodels.User) (internalmodels.User, internalmodels.StatusCode)
	GetUserForItsNickname(user internalmodels.User) (internalmodels.User, internalmodels.StatusCode)
	DeactivateUser(user internalmodels.User) internalmodels.StatusCode
	UpdateUser(user internalmodels.User) (internalmodels.User, internalmodels.StatusCode)
	FindUsers(nickname string) ([]internalmodels.User, internalmodels.StatusCode)
}

type UserUsecase interface {
	SignUp(request models.SignUpPersonRequest) (internalmodels.User, internalmodels.StatusCode)
	Login(request models.LoginRequest) (internalmodels.User, internalmodels.StatusCode)
	GetUser(request models.GetUserRequest) (internalmodels.User, internalmodels.StatusCode)
	DeactivateUser(user internalmodels.User) internalmodels.StatusCode
	UpdateUser(userToUpdate internalmodels.User,
		request models.UpdateUserRequest) (internalmodels.User, internalmodels.StatusCode)
	FindUsers(nickname string) ([]internalmodels.User, internalmodels.StatusCode)
}

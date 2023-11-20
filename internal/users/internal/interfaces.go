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
	ActivateUser(user internalmodels.User) internalmodels.StatusCode
}

type UserUsecase interface {
	SignUp(request models.SignUpPersonRequest) (internalmodels.Session, internalmodels.StatusCode)
	Login(request models.LoginRequest) (internalmodels.Session, internalmodels.StatusCode)
	Logout(session internalmodels.Session) internalmodels.StatusCode
	GetUser(request models.GetUserRequest) (internalmodels.User, internalmodels.StatusCode)
	DeactivateUser(user internalmodels.User) internalmodels.StatusCode
	UpdateUser(userToUpdate internalmodels.User,
		request models.UpdateUserRequest) (internalmodels.User, internalmodels.StatusCode)
	FindUsers(nickname string) ([]internalmodels.User, internalmodels.StatusCode)
	ActivateUser(session internalmodels.Session, code string) internalmodels.StatusCode
	GetSession(session internalmodels.Session) (internalmodels.Session, internalmodels.StatusCode)
}

type SessionRepo interface {
	CreateSession(user internalmodels.User,
		sessionType string) (internalmodels.Session, internalmodels.StatusCode)
	GetSession(session internalmodels.Session) (internalmodels.Session, internalmodels.StatusCode)
	DeleteSession(session internalmodels.Session) internalmodels.StatusCode
}

type ActivationRepo interface {
	CreateActivationCode(session internalmodels.Session) (string, internalmodels.StatusCode)
	CheckActivationCode(session internalmodels.Session,
		activationCode string) (bool, internalmodels.StatusCode)
}

type MailerRepo interface {
	PutActivationTask(request internalmodels.ActivationTask) internalmodels.StatusCode
}

package usecase

import (
	"our-little-chatik/internal/gateway/internal/delivery"
	"our-little-chatik/internal/models"
)

type GatewayUsecase struct {
	auth     delivery.AuthHandler
	userData delivery.UserDataHandler
}

func NewGatewayUsecasse(auth delivery.AuthHandler,
	userData delivery.UserDataHandler) *GatewayUsecase {
	return &GatewayUsecase{
		auth:     auth,
		userData: userData,
	}
}

func (u GatewayUsecase) SignUp(user *models.User) (*models.Session, error) {
	newUser, err := u.userData.AddUser(user)
	if err != nil {
		return nil, err
	}

	session, err := u.auth.AddUser(*newUser)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u GatewayUsecase) SignIn(user *models.User) (*models.Session, error) {
	err := u.userData.CheckUser(user)
	if err != nil {
		return nil, err
	}

	session, err := u.auth.AddUser(*user)
	if err != nil {
		return nil, err
	}
	return session, err
}

func (u GatewayUsecase) GetSessionFromUser(user models.User) (*models.Session, error) {
	session, err := u.auth.GetSession(user)
	if err != nil {
		return nil, err
	}
	return session, err
}

func (u GatewayUsecase) GetUserFromSession(session models.Session) (*models.User, error) {
	user, err := u.auth.GetUser(session)
	if err != nil {
		return nil, err
	}
	return user, nil
}

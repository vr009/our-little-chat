package usecase

import (
	"our-little-chatik/internal/gateway/internal/delivery"
	"our-little-chatik/internal/models"

	"github.com/golang/glog"
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
		glog.Warningf("Failed to add a user in user-data service. Error message: %s.", err.Error())
		return nil, err
	}

	session, err := u.auth.AddUser(*newUser)
	if err != nil {
		glog.Warningf("Failed to add a user in auth-data service. Error message: %s.", err.Error())
		return nil, err
	}

	return session, nil
}

func (u GatewayUsecase) SignIn(user *models.User) (*models.Session, error) {
	err := u.userData.CheckUser(user)
	if err != nil {
		glog.Errorf("Failed to check a user in user-data service. Error message: %s.", err.Error())
		return nil, err
	}

	session, err := u.auth.AddUser(*user)
	if err != nil {
		glog.Errorf("Failed to add a user in user-data service. Error message: %s.", err.Error())
		return nil, err
	}
	return session, err
}

func (u GatewayUsecase) LogOut(session models.Session) error {
	err := u.auth.RemoveUser(session)
	if err != nil {
		glog.Errorf("Failed to add a user in user-data service. Error message: %s.", err.Error())
		return err
	}
	return nil
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

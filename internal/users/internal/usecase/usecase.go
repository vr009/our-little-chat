package usecase

import (
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/models"
	models2 "our-little-chatik/internal/users/internal/models"
	"time"

	"github.com/google/uuid"
	"our-little-chatik/internal/users/internal"
)

type UserUsecase struct {
	userRepo       internal.UserRepo
	sessionRepo    internal.SessionRepo
	activationRepo internal.ActivationRepo
	mailerRepo     internal.MailerRepo
}

func NewUserUsecase(userRepo internal.UserRepo,
	sessionRepo internal.SessionRepo, activationRepo internal.ActivationRepo,
	mailerRepo internal.MailerRepo) *UserUsecase {
	return &UserUsecase{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		activationRepo: activationRepo,
		mailerRepo:     mailerRepo,
	}
}

// TODO add degradation handling
func (uc *UserUsecase) SignUp(request models2.SignUpPersonRequest) (models.Session, models.StatusCode) {
	user := models.User{
		Name:      *request.Name,
		Nickname:  *request.Nickname,
		Surname:   *request.Surname,
		Avatar:    *request.Avatar,
		Activated: false,
	}
	err := user.Password.Set(*request.Password)
	if err != nil {
		return models.Session{}, models.InternalError
	}
	user.ID = uuid.New()
	user.RegisteredAt = time.Now()

	user, status := uc.userRepo.CreateUser(user)
	if status != models.OK {
		return models.Session{}, status
	}

	session, status := uc.sessionRepo.CreateSession(user, models.ActivationSession)
	if status != models.OK {
		return models.Session{}, status
	}

	code, status := uc.activationRepo.CreateActivationCode(session)
	if status != models.OK {
		return models.Session{}, models.InternalError
	}

	slog.Info("code %s", code)
	uc.mailerRepo.PutActivationTask(models.ActivationTask{
		ActivationCode: code,
		Receiver:       user.Email,
	})
	return session, status
}

func (uc *UserUsecase) Login(request models2.LoginRequest) (models.Session, models.StatusCode) {
	user, status := uc.userRepo.GetUserForItsNickname(models.User{Nickname: *request.Nickname})
	if status != models.OK {
		return models.Session{}, status
	}

	if !user.Activated {
		return models.Session{}, models.InActivated
	}

	ok, err := user.Password.Matches(*request.Password)
	if err != nil {
		return models.Session{}, models.InternalError
	}
	if !ok {
		return models.Session{}, models.Unauthorized
	}

	return uc.sessionRepo.CreateSession(user, models.PlainSession)
}

func (uc *UserUsecase) Logout(session models.Session) models.StatusCode {
	return uc.sessionRepo.DeleteSession(session)
}

func (uc *UserUsecase) ActivateUser(session models.Session,
	code string) models.StatusCode {
	if ok, status := uc.activationRepo.CheckActivationCode(session, code); status != models.OK || !ok {
		return models.InActivated
	}
	s, status := uc.sessionRepo.GetSession(session)
	if status != models.OK {
		return status
	}
	return uc.userRepo.ActivateUser(models.User{ID: s.UserID})
}

func (uc *UserUsecase) GetSession(session models.Session) (models.Session, models.StatusCode) {
	return uc.sessionRepo.GetSession(session)
}

func (uc *UserUsecase) GetUser(request models2.GetUserRequest) (models.User, models.StatusCode) {
	return uc.userRepo.GetUserForItsID(models.User{ID: request.UserID})
}

func (uc *UserUsecase) DeactivateUser(user models.User) models.StatusCode {
	return uc.userRepo.DeactivateUser(user)
}

func (uc *UserUsecase) UpdateUser(userToUpdate models.User,
	request models2.UpdateUserRequest) (models.User, models.StatusCode) {
	oldUser, status := uc.userRepo.GetUserForItsID(userToUpdate)
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
	return uc.userRepo.UpdateUser(newUser)
}

func (uc *UserUsecase) FindUsers(name string) ([]models.User, models.StatusCode) {
	return uc.userRepo.FindUsers(name)
}

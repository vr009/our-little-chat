package delivery

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg/proto/session"
	"our-little-chatik/internal/pkg/proto/users"
	"our-little-chatik/internal/users/internal"
	"our-little-chatik/internal/users/internal/models"
)

type UserGRPCHandler struct {
	useCase internal.UserUsecase
	users.UnimplementedUsersServer
	session.UnimplementedSessionServer
}

func NewUserGRPCHandler(useCase internal.UserUsecase) *UserGRPCHandler {
	return &UserGRPCHandler{
		useCase: useCase,
	}
}

func (h UserGRPCHandler) GetUser(ctx context.Context, request *users.GetUserRequest) (*users.UserResponse, error) {
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		return nil, err
	}
	user, status := h.useCase.GetUser(models.GetUserRequest{UserID: userID})
	if status != models2.OK {
		return nil, fmt.Errorf("failed to get user")
	}
	return &users.UserResponse{
		UserID:    user.ID.String(),
		Name:      user.Name,
		Surname:   user.Surname,
		Nickname:  user.Nickname,
		Activated: user.Activated,
		Avatar:    user.Avatar,
	}, nil
}

func (h UserGRPCHandler) GetSession(ctx context.Context,
	request *session.GetSessionRequest) (*session.SessionResponse, error) {
	sessionID, err := uuid.Parse(request.SessionID)
	if err != nil {
		return nil, err
	}
	s, status := h.useCase.GetSession(models2.Session{ID: sessionID})
	if status != models2.OK {
		return nil, fmt.Errorf("failed to get user")
	}
	return &session.SessionResponse{
		SessionID: s.ID.String(),
		UserID:    s.UserID.String(),
		Type:      s.Type,
	}, nil
}

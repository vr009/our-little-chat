package delivery

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg/proto/users"
	"our-little-chatik/internal/users/internal"
	"our-little-chatik/internal/users/internal/models"
)

type UserGRPCHandler struct {
	useCase internal.UserUsecase
	users.UnimplementedUsersServer
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

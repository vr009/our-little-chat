package repo

import (
	"context"
	"github.com/google/uuid"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg/proto/users"
)

type UserDataClient struct {
	cl users.UsersClient
}

func NewUserDataClient(cl users.UsersClient) *UserDataClient {
	return &UserDataClient{
		cl: cl,
	}
}

func (c UserDataClient) GetUser(user models.User) (models.User, models.StatusCode) {
	resp, err := c.cl.GetUser(context.Background(),
		&users.GetUserRequest{UserID: user.ID.String()})
	if err != nil {
		return models.User{}, models.NotFound
	}
	user = models.User{
		Name:      resp.Name,
		Nickname:  resp.Nickname,
		Surname:   resp.Surname,
		Avatar:    resp.Avatar,
		Activated: resp.Activated,
	}
	user.ID, err = uuid.Parse(resp.UserID)
	if err != nil {
		return models.User{}, models.InternalError
	}
	return user, models.OK
}

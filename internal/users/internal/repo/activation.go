package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"our-little-chatik/internal/models"
	"time"
)

type ActivationRepo struct {
	cl *redis.Client
}

func NewActivationRepo(cl *redis.Client) *ActivationRepo {
	return &ActivationRepo{
		cl: cl,
	}
}

func (ar *ActivationRepo) CreateActivationCode(session models.Session) (string, models.StatusCode) {
	code := uuid.New().String()
	err := ar.cl.Set(context.Background(), session.ID.String(), code, time.Hour*3).Err()
	if err != nil {
		return "", models.InternalError
	}
	return code, models.OK
}

func (ar *ActivationRepo) CheckActivationCode(session models.Session,
	activationCode string) (bool, models.StatusCode) {
	storedCode, err := ar.cl.Get(context.Background(), session.ID.String()).Result()
	if err != nil {
		return false, models.InternalError
	}
	return storedCode == activationCode, models.OK
}

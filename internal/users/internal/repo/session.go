package repo

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"our-little-chatik/internal/models"
	"time"
)

type SessionRepo struct {
	cl *redis.Client
}

func NewSessionRepo(cl *redis.Client) *SessionRepo {
	return &SessionRepo{
		cl: cl,
	}
}

func (sr *SessionRepo) CreateSession(user models.User,
	sessionType string) (models.Session, models.StatusCode) {
	newSession := models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		Type:      sessionType,
	}
	body, _ := json.Marshal(newSession)
	err := sr.cl.Set(context.Background(), newSession.ID.String(), string(body), time.Hour*48)
	if err != nil {
		return models.Session{}, models.InternalError
	}
	return newSession, models.OK
}

func (sr *SessionRepo) GetSession(session models.Session) (models.Session, models.StatusCode) {
	val, err := sr.cl.Get(context.Background(), session.ID.String()).Result()
	if err != nil {
		return models.Session{}, models.NotFound
	}
	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		return models.Session{}, models.InActivated
	}
	return session, models.OK
}

func (sr *SessionRepo) DeleteSession(session models.Session) models.StatusCode {
	_, err := sr.cl.Del(context.Background(), session.ID.String()).Result()
	if err != nil {
		return models.NotFound
	}
	return models.OK
}

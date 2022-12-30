package repo

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"time"

	"our-little-chatik/internal/auth/internal/models"
	models2 "our-little-chatik/internal/models"

	"github.com/go-redis/redis/v9"
	"github.com/golang/glog"
	"github.com/google/uuid"
)

type DataBase struct {
	Client *redis.Client
	TTL    int
}

func MD5(data string) string {
	h := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", h)
}

func NewDataBase(Client *redis.Client, TTL int) *DataBase {
	return &DataBase{
		Client: Client,
		TTL:    TTL,
	}
}

func (db *DataBase) CreateSession(session models2.Session) (models2.Session, models.StatusCode) {
	token := MD5(session.UserID.String())

	session.Token = token

	err := db.Client.Set(context.Background(), session.UserID.String(), session.Token, time.Minute*time.Duration(db.TTL)).Err()

	if err != nil {
		glog.Errorf("failed to add a session: %s", err.Error())
		return models2.Session{}, models.Conflict
	}

	err = db.Client.Set(context.Background(), session.Token, session.UserID.String(), time.Minute*time.Duration(db.TTL)).Err()

	if err != nil {
		glog.Errorf("failed to add a session: %s", err.Error())
		return models2.Session{}, models.Conflict
	}

	return session, models.OK
}

func (db *DataBase) GetToken(session models2.Session) (models2.Session, models.StatusCode) {

	cmdToken := db.Client.Get(context.Background(), session.UserID.String())
	if cmdToken.Err() != nil {
		return models2.Session{}, models.NotFound
	}

	token, err := cmdToken.Result()
	if err != nil {
		return models2.Session{}, models.InternalError
	}

	session.Token = token

	return session, models.OK
}

func (db *DataBase) GetUser(session models2.Session) (models2.Session, models.StatusCode) {

	cmdUser := db.Client.Get(context.Background(), session.Token)

	if cmdUser.Err() != nil {
		return models2.Session{}, models.NotFound
	}

	userId, err := cmdUser.Result()

	if err != nil {
		log.Print("error of cmd.Result parsing")
		return models2.Session{}, models.InternalError
	}

	session.UserID, err = uuid.Parse(userId)

	if err != nil {
		log.Print("error of UUID parsing")
		return models2.Session{}, models.InternalError
	}

	return session, models.OK

}

func (db *DataBase) DeleteSession(session models2.Session) models.StatusCode {

	s := models2.Session{}

	cmd := db.Client.Get(context.Background(), session.Token)

	if cmd.Err() != nil {
		return models.NotFound
	}

	value, err := cmd.Result()

	if err != nil {
		log.Print("error of cmd.Result parsing")
		return models.InternalError
	}

	userIdFromUuid, err := uuid.Parse(value)

	if err != nil {
		log.Print("error of UUID parsing")
		return models.InternalError
	}

	s.Token = session.Token
	s.UserID = userIdFromUuid

	err = db.Client.Del(context.Background(), s.Token).Err()

	if err != nil {
		log.Print("error of deleting UserID")
		return models.InternalError
	}

	err = db.Client.Del(context.Background(), s.UserID.String()).Err()
	if err != nil {
		log.Print("error of deleting Token")
		return models.InternalError
	}

	return models.OK
}

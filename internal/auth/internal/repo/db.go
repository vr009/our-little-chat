package repo

import (
	"auth/internal/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"log"
)

type DataBase struct {
	Client *redis.Client
}

//redis-server

func NewDataBase(Client *redis.Client) *DataBase {
	return &DataBase{
		Client: Client,
	}
}

func (db *DataBase) CreateSession(session models.Session) (models.Session, models.StatusCode) {
	token, err := uuid.NewUUID()

	if err != nil {
		log.Print("error with generation of UUID")
		return models.Session{}, models.InternalError
	}

	session.Token = token.String()

	err = db.Client.Set(context.Background(), session.UserID.String(), session.Token, 0).Err()

	if err != nil {
		return models.Session{}, models.Conflict
	}

	err = db.Client.Set(context.Background(), session.Token, session.UserID.String(), 0).Err()

	if err != nil {
		return models.Session{}, models.Conflict
	}

	return session, models.OK
}

func (db *DataBase) GetToken(session models.Session) (models.Session, models.StatusCode) {

	fmt.Println(session.UserID.String())

	cmd := db.Client.Get(context.Background(), session.UserID.String())

	if cmd.Err() != nil {
		return models.Session{}, models.NotFound
	}
	value, err := cmd.Result()
	if err != nil {
		return models.Session{}, models.InternalError
	}

	return models.Session{
		Token:  value,
		UserID: session.UserID,
	}, models.OK

}

func (db *DataBase) GetUser(session models.Session) (models.Session, models.StatusCode) {

	fmt.Println(session.Token)

	cmd := db.Client.Get(context.Background(), session.Token)

	if cmd.Err() != nil {
		return models.Session{}, models.NotFound
	}

	value, err := cmd.Result()

	if err != nil {
		log.Print("error of cmd.Result parsing")
		return models.Session{}, models.InternalError
	}

	uuidFormString, err := uuid.Parse(value)

	if err != nil {
		log.Print("error of UUID parsing")
		return models.Session{}, models.InternalError
	}

	fmt.Println(value)
	fmt.Println(uuidFormString)

	return models.Session{
		Token:  session.Token,
		UserID: uuidFormString,
	}, models.OK
}

func (db *DataBase) DeleteSession(session models.Session) models.StatusCode {

	s := models.Session{}

	fmt.Println(session.Token)

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

package repo

import (
	"user_data/internal/models"
)

type DataBase struct {
	Client string
	TTL    int
}

//todo Поправить Client –> заменяем на PostgresSQL

func NewDataBase(Client string, TTL int) *DataBase {
	return &DataBase{
		Client: Client,
		TTL:    TTL,
	}
}

func (db *DataBase) GetAllUsers(userData models.UserData) (models.UserData, models.StatusCode) {
	return userData, models.OK
}

func (db *DataBase) GetUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return userData, models.OK
}

func (db *DataBase) CreateUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return userData, models.OK
}

func (db *DataBase) DeleteUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return userData, models.OK
}

func (db *DataBase) UpdateUser(userData models.UserData) (models.UserData, models.StatusCode) {
	return userData, models.OK
}

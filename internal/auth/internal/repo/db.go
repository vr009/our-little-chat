package repo

import "auth/internal/models"

type DataBase struct {
}

func NewDataBase() *DataBase {
	return &DataBase{}
}

func (db *DataBase) CreateSession(session models.Session) models.Session {
	return models.Session{}
}

func (db *DataBase) DeleteSession(session models.Session) {
	return
}

func (db *DataBase) GetSession(session models.Session) models.Session {
	return models.Session{}
}

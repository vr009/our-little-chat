package repo

import (
	"context"
	"database/sql"
	"log"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/user_data/internal"

	"golang.org/x/exp/slog"
)

const (
	InsertQuery = "INSERT INTO users(user_id, nickname, name, surname, password, last_auth, registered, avatar) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8);"
	DeleteQuery    = "DELETE FROM users WHERE user_id=$1;"
	UpdateQuery    = "UPDATE users SET nickname=$1, name=$2, surname=$3, avatar=$4 WHERE user_id=$5;"
	GetQuery       = "SELECT user_id, nickname, name, surname, last_auth, registered, avatar  FROM users WHERE user_id=$1;"
	GetNameQuery   = "SELECT user_id, nickname, name, surname, password, last_auth, registered, avatar  FROM users WHERE nickname=$1;"
	ListQuery      = "SELECT user_id, nickname, avatar FROM users;"
	FindUsersQuery = "SELECT user_id, nickname, name, surname, avatar FROM users WHERE nickname LIKE LOWER($1 || '%') LIMIT 10"
)

type PersonRepo struct {
	pool internal.DB
}

func NewPersonRepo(pool internal.DB) *PersonRepo {
	return &PersonRepo{
		pool: pool,
	}
}

func (pr *PersonRepo) CreateUser(person models2.UserData) (models2.UserData, models2.StatusCode) {
	slog.Info("Creating user", "user_data data", slog.AnyValue(person))
	_, err := pr.pool.Exec(context.Background(),
		InsertQuery,
		person.UserID,
		person.Nickname,
		person.Name,
		person.Surname,
		person.Password,
		person.LastAuth,
		person.Registered,
		person.Avatar,
	)
	if err != nil {
		slog.Error(err.Error())
		return models2.UserData{}, models2.BadRequest
	}
	person.Password = ""
	return person, models2.OK
}

func (pr *PersonRepo) DeleteUser(person models2.UserData) models2.StatusCode {
	_, err := pr.pool.Exec(context.Background(), DeleteQuery, person.UserID)
	if err != nil {
		return models2.InternalError
	}
	return models2.Deleted
}

func (pr *PersonRepo) UpdateUser(personNew models2.UserData) (models2.UserData, models2.StatusCode) {
	_, err := pr.pool.Exec(context.Background(), UpdateQuery,
		personNew.Nickname, personNew.Name, personNew.Surname, personNew.Avatar,
		personNew.UserID)
	if err != nil {
		return models2.UserData{}, models2.BadRequest
	}
	return personNew, models2.OK
}

func (pr *PersonRepo) GetUser(person models2.UserData) (models2.UserData, models2.StatusCode) {
	log.Println("SEARCHING USER ID !!!!!!!!!!!!", person.UserID)
	rows := pr.pool.QueryRow(context.Background(), GetQuery, person.UserID)
	err := rows.Scan(&person.UserID, &person.Nickname,
		&person.Name, &person.Surname, &person.LastAuth,
		&person.Registered, &person.Avatar)
	if err != nil {
		slog.Error("user_data not found: " + err.Error())
		return models2.UserData{}, models2.NotFound
	}
	log.Println("FOUND USER ID !!!!!!!!!!!!", person.UserID)
	return person, models2.OK
}

func (pr *PersonRepo) GetUserForItsName(person models2.UserData) (models2.UserData, models2.StatusCode) {
	rows := pr.pool.QueryRow(context.Background(), GetNameQuery, person.Nickname)
	err := rows.Scan(&person.UserID, &person.Nickname,
		&person.Name, &person.Surname, &person.Password, &person.LastAuth,
		&person.Registered, &person.Avatar)
	if err != nil {
		slog.Error("user_data not found: " + err.Error())
		return models2.UserData{}, models2.NotFound
	}
	return person, models2.OK
}

func (pr *PersonRepo) GetAllUsers() ([]models2.UserData, models2.StatusCode) {
	rows, err := pr.pool.Query(context.Background(), ListQuery)
	if err != nil && err != sql.ErrNoRows {
		return nil, models2.InternalError
	}
	defer rows.Close()
	list := make([]models2.UserData, 0)
	for rows.Next() {
		person := models2.UserData{}
		rows.Scan(&person.UserID, &person.Nickname, &person.Avatar)
		list = append(list, person)
	}
	return list, models2.OK
}

func (pr *PersonRepo) FindUser(name string) ([]models2.UserData, models2.StatusCode) {
	rows, err := pr.pool.Query(context.Background(), FindUsersQuery, name)
	if err != nil && err != sql.ErrNoRows {
		slog.Error(err.Error())
		return nil, models2.InternalError
	}
	defer rows.Close()
	list := make([]models2.UserData, 0)
	for rows.Next() {
		person := models2.UserData{}
		rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.Avatar)
		list = append(list, person)
	}
	slog.Info("List:", "users", slog.AnyValue(list))
	return list, models2.OK
}

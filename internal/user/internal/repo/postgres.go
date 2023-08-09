package repo

import (
	"context"
	"database/sql"
	"our-little-chatik/internal/user/internal"

	"our-little-chatik/internal/user/internal/models"

	"golang.org/x/exp/slog"
)

const (
	InsertQuery = "INSERT INTO users(user_id, nickname, name, surname, password, last_auth, registered, avatar) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8);"
	DeleteQuery    = "DELETE FROM users WHERE user_id=$1;"
	UpdateQuery    = "UPDATE users SET nickname=$1, name=$2, surname=$3, avatar=$4 WHERE user_id=$5;"
	GetQuery       = "SELECT user_id, nickname, name, surname, last_auth, registered, avatar  FROM users WHERE user_id=$1;"
	GetNameQuery   = "SELECT user_id, nickname, name, surname, last_auth, registered, avatar  FROM users WHERE nickname=$1;"
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

func (pr *PersonRepo) CreateUser(person models.UserData) (models.UserData, models.StatusCode) {
	slog.Info("Creating user", "user data", slog.AnyValue(person))
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
		return models.UserData{}, models.BadRequest
	}
	person.Password = ""
	return person, models.OK
}

func (pr *PersonRepo) DeleteUser(person models.UserData) models.StatusCode {
	_, err := pr.pool.Exec(context.Background(), DeleteQuery, person.UserID)
	if err != nil {
		return models.InternalError
	}
	return models.Deleted
}

func (pr *PersonRepo) UpdateUser(personNew models.UserData) (models.UserData, models.StatusCode) {
	_, err := pr.pool.Exec(context.Background(), UpdateQuery,
		personNew.Nickname, personNew.Name, personNew.Surname, personNew.Avatar,
		personNew.UserID)
	if err != nil {
		return models.UserData{}, models.BadRequest
	}
	return personNew, models.OK
}

func (pr *PersonRepo) GetUser(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.pool.QueryRow(context.Background(), GetQuery, person.UserID)
	slog.Info("SEARCHING FOR " + person.UserID.String())
	err := rows.Scan(&person.UserID, &person.Nickname,
		&person.Name, &person.Surname, &person.LastAuth,
		&person.Registered, &person.Avatar)
	if err != nil {
		slog.Error("user not found: " + err.Error())
		return models.UserData{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetUserForItsName(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.pool.QueryRow(context.Background(), GetNameQuery, person.Nickname)
	err := rows.Scan(&person.UserID, &person.Nickname,
		&person.Name, &person.Surname, &person.LastAuth,
		&person.Registered, &person.Avatar)
	if err != nil {
		slog.Error("user not found: " + err.Error())
		return models.UserData{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetAllUsers() ([]models.UserData, models.StatusCode) {
	rows, err := pr.pool.Query(context.Background(), ListQuery)
	if err != nil && err != sql.ErrNoRows {
		return nil, models.InternalError
	}
	defer rows.Close()
	list := make([]models.UserData, 0)
	for rows.Next() {
		person := models.UserData{}
		rows.Scan(&person.UserID, &person.Nickname, &person.Avatar)
		list = append(list, person)
	}
	return list, models.OK
}

func (pr *PersonRepo) FindUser(name string) ([]models.UserData, models.StatusCode) {
	rows, err := pr.pool.Query(context.Background(), FindUsersQuery, name)
	if err != nil && err != sql.ErrNoRows {
		slog.Error(err.Error())
		return nil, models.InternalError
	}
	defer rows.Close()
	list := make([]models.UserData, 0)
	for rows.Next() {
		person := models.UserData{}
		rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.Avatar)
		list = append(list, person)
	}
	slog.Info("List:", "users", slog.AnyValue(list))
	return list, models.OK
}

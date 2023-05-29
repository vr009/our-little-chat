package repo

import (
	"context"
	"database/sql"

	"our-little-chatik/internal/user_data/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

const (
	InsertQuery = "INSERT INTO users(user_id, nickname, name, surname, password, last_auth, registered, avatar, contact_list) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);"
	DeleteQuery    = "DELETE FROM users WHERE user_id=$1;"
	UpdateQuery    = "UPDATE users SET nickname=$1, name=$2, surname=$3, last_auth=$4, registered=$5, avatar=$6, contact_list=$7 WHERE user_id=$8;"
	GetQuery       = "SELECT user_id, nickname, name, surname, last_auth, registered, avatar, contact_list  FROM users WHERE user_id=$1;"
	GetNameQuery   = "SELECT user_id, nickname, name, surname, password  FROM users WHERE nickname=$1;"
	ListQuery      = "SELECT * FROM users;"
	FindUsersQuery = "SELECT user_id, nickname, name, surname FROM users WHERE nickname LIKE LOWER($1 || '%') LIMIT 10"
)

type PersonRepo struct {
	pool *pgxpool.Pool
}

func NewPersonRepo(pool *pgxpool.Pool) *PersonRepo {
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
		person.ContactList,
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
	return models.OK
}

func (pr *PersonRepo) UpdateUser(personNew models.UserData) (models.UserData, models.StatusCode) {
	personOld, status := pr.GetUser(personNew)
	if status != models.OK {
		return personOld, models.NotFound
	}

	if personOld.Nickname != personNew.Nickname && personNew.Nickname != "" {
		personOld.Nickname = personNew.Nickname
	}

	if personOld.LastAuth != personNew.LastAuth {
		personOld.LastAuth = personNew.LastAuth
	}

	_, err := pr.pool.Exec(context.Background(), UpdateQuery, personNew.Nickname, &personNew.Name, &personNew.Surname, personNew.LastAuth, personNew.Registered, personNew.Avatar, personNew.ContactList, personNew.UserID)
	if err != nil {
		return personOld, models.BadRequest
	}
	personNew = personOld
	return personNew, models.OK
}

func (pr *PersonRepo) GetUser(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.pool.QueryRow(context.Background(), GetQuery, person.UserID)
	slog.Info("SEARCHING FOR " + person.UserID.String())
	err := rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
	if err != nil {
		slog.Error("user not found: " + err.Error())
		return models.UserData{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetUserForItsName(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.pool.QueryRow(context.Background(), GetNameQuery, person.Nickname)
	err := rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.Password)
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
		rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
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
		rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname)
		list = append(list, person)
	}
	slog.Info("List:", "users", slog.AnyValue(list))
	return list, models.OK
}

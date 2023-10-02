package repo

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/exp/slog"
	models2 "our-little-chatik/internal/models"
)

const (
	InsertQuery = "INSERT INTO users(user_id, nickname, user_name, surname, password, avatar) " +
		"VALUES($1, $2, $3, $4, $5, $6) RETURNING registered;"
	DeleteQuery    = "UPDATE users SET activated=false WHERE user_id=$1;"
	UpdateQuery    = "UPDATE users SET nickname=$1, user_name=$2, surname=$3, avatar=$4, password=$5 WHERE user_id=$6;"
	GetQuery       = "SELECT user_id, nickname, user_name, surname, password, registered, avatar  FROM users WHERE user_id=$1;"
	GetNameQuery   = "SELECT user_id, nickname, user_name, surname, password, registered, avatar  FROM users WHERE nickname=$1;"
	ListQuery      = "SELECT user_id, nickname, avatar FROM users;"
	FindUsersQuery = "SELECT user_id, nickname, user_name, surname, avatar FROM users WHERE nickname LIKE LOWER($1 || '%') LIMIT 10"
)

type PersonRepo struct {
	pool *sql.DB
}

func NewPersonRepo(pool *sql.DB) *PersonRepo {
	return &PersonRepo{
		pool: pool,
	}
}

func (pr *PersonRepo) CreateUser(person models2.User) (models2.User, models2.StatusCode) {
	slog.Info("Creating user", "user_data data", slog.AnyValue(person))
	err := pr.pool.QueryRowContext(context.Background(),
		InsertQuery,
		person.ID,
		person.Nickname,
		person.Name,
		person.Surname,
		person.Password.Hash,
		person.Avatar,
	).Scan(&person.Registered)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return models2.User{}, models2.Conflict
		default:
			return models2.User{}, models2.InternalError
		}
	}
	return person, models2.OK
}

func (pr *PersonRepo) DeleteUser(person models2.User) models2.StatusCode {
	_, err := pr.pool.ExecContext(context.Background(), DeleteQuery, person.ID)
	if err != nil {
		return models2.InternalError
	}
	return models2.Deleted
}

func (pr *PersonRepo) UpdateUser(personNew models2.User) (models2.User, models2.StatusCode) {
	_, err := pr.pool.ExecContext(context.Background(), UpdateQuery,
		personNew.Nickname, personNew.Name, personNew.Surname, personNew.Avatar,
		personNew.Password.Hash, personNew.ID)
	if err != nil {
		return models2.User{}, models2.BadRequest
	}
	return personNew, models2.OK
}

func (pr *PersonRepo) GetUser(person models2.User) (models2.User, models2.StatusCode) {
	rows := pr.pool.QueryRowContext(context.Background(), GetQuery, person.ID)
	err := rows.Scan(&person.ID, &person.Nickname,
		&person.Name, &person.Surname, &person.Password.Hash,
		&person.Registered, &person.Avatar)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models2.User{}, models2.NotFound
		default:
			return models2.User{}, models2.InternalError
		}
	}
	return person, models2.OK
}

func (pr *PersonRepo) GetUserForItsName(person models2.User) (models2.User, models2.StatusCode) {
	rows := pr.pool.QueryRowContext(context.Background(), GetNameQuery, person.Nickname)
	err := rows.Scan(&person.ID, &person.Nickname,
		&person.Name, &person.Surname, &person.Password.Hash,
		&person.Registered, &person.Avatar)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models2.User{}, models2.NotFound
		default:
			return models2.User{}, models2.InternalError
		}
	}
	return person, models2.OK
}

func (pr *PersonRepo) GetAllUsers() ([]models2.User, models2.StatusCode) {
	rows, err := pr.pool.QueryContext(context.Background(), ListQuery)
	if err != nil && err != sql.ErrNoRows {
		return nil, models2.InternalError
	}
	defer rows.Close()
	list := make([]models2.User, 0)
	for rows.Next() {
		person := models2.User{}
		rows.Scan(&person.ID, &person.Nickname, &person.Avatar)
		list = append(list, person)
	}
	return list, models2.OK
}

func (pr *PersonRepo) FindUser(name string) ([]models2.User, models2.StatusCode) {
	rows, err := pr.pool.QueryContext(context.Background(), FindUsersQuery, name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, models2.NotFound
		default:
			return nil, models2.InternalError
		}
	}
	defer rows.Close()
	list := make([]models2.User, 0)
	for rows.Next() {
		person := models2.User{}
		rows.Scan(&person.ID, &person.Nickname, &person.Name, &person.Surname, &person.Avatar)
		list = append(list, person)
	}
	slog.Info("List:", "users", slog.AnyValue(list))
	return list, models2.OK
}

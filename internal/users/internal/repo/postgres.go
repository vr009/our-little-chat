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
	FindUsersQuery = "SELECT user_id, nickname, user_name, surname, avatar FROM users WHERE nickname LIKE LOWER($1 || '%') LIMIT 10"
)

type UserRepo struct {
	pool *sql.DB
}

func NewUserRepo(pool *sql.DB) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (pr *UserRepo) CreateUser(user models2.User) (models2.User, models2.StatusCode) {
	slog.Info("Creating user", "user_data data", slog.AnyValue(user))
	err := pr.pool.QueryRowContext(context.Background(),
		InsertQuery,
		user.ID,
		user.Nickname,
		user.Name,
		user.Surname,
		user.Password.Hash,
		user.Avatar,
	).Scan(&user.Registered)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return models2.User{}, models2.Conflict
		default:
			return models2.User{}, models2.InternalError
		}
	}
	return user, models2.OK
}

func (pr *UserRepo) DeactivateUser(user models2.User) models2.StatusCode {
	_, err := pr.pool.ExecContext(context.Background(), DeleteQuery, user.ID)
	if err != nil {
		return models2.InternalError
	}
	return models2.Deleted
}

func (pr *UserRepo) UpdateUser(userNew models2.User) (models2.User, models2.StatusCode) {
	_, err := pr.pool.ExecContext(context.Background(), UpdateQuery,
		userNew.Nickname, userNew.Name, userNew.Surname, userNew.Avatar,
		userNew.Password.Hash, userNew.ID)
	if err != nil {
		return models2.User{}, models2.BadRequest
	}
	return userNew, models2.OK
}

func (pr *UserRepo) GetUserForItsID(user models2.User) (models2.User, models2.StatusCode) {
	rows := pr.pool.QueryRowContext(context.Background(), GetQuery, user.ID)
	err := rows.Scan(&user.ID, &user.Nickname,
		&user.Name, &user.Surname, &user.Password.Hash,
		&user.Registered, &user.Avatar)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models2.User{}, models2.NotFound
		default:
			return models2.User{}, models2.InternalError
		}
	}
	return user, models2.OK
}

func (pr *UserRepo) GetUserForItsNickname(user models2.User) (models2.User, models2.StatusCode) {
	rows := pr.pool.QueryRowContext(context.Background(), GetNameQuery, user.Nickname)
	err := rows.Scan(&user.ID, &user.Nickname,
		&user.Name, &user.Surname, &user.Password.Hash,
		&user.Registered, &user.Avatar)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models2.User{}, models2.NotFound
		default:
			return models2.User{}, models2.InternalError
		}
	}
	return user, models2.OK
}

func (pr *UserRepo) FindUsers(name string) ([]models2.User, models2.StatusCode) {
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
		user := models2.User{}
		rows.Scan(&user.ID, &user.Nickname, &user.Name, &user.Surname, &user.Avatar)
		list = append(list, user)
	}
	slog.Info("List:", "users", slog.AnyValue(list))
	return list, models2.OK
}

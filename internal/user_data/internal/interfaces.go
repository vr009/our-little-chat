package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"our-little-chatik/internal/user_data/internal/models"
)

type UserdataRepo interface {
	GetAllUsers() ([]models.UserData, models.StatusCode)
	CreateUser(userData models.UserData) (models.UserData, models.StatusCode)
	GetUser(userData models.UserData) (models.UserData, models.StatusCode)
	GetUserForItsName(userData models.UserData) (models.UserData, models.StatusCode)
	DeleteUser(userData models.UserData) models.StatusCode
	UpdateUser(personNew models.UserData) (models.UserData, models.StatusCode)
	FindUser(name string) ([]models.UserData, models.StatusCode)
}

type UserdataUseCase interface {
	GetAllUsers() ([]models.UserData, models.StatusCode)
	CreateUser(userData models.UserData) (models.UserData, models.StatusCode)
	GetUser(userData models.UserData) (models.UserData, models.StatusCode)
	DeleteUser(userData models.UserData) models.StatusCode
	UpdateUser(userData models.UserData) (models.UserData, models.StatusCode)
	CheckUser(userData models.UserData) (models.UserData, models.StatusCode)
	FindUser(name string) ([]models.UserData, models.StatusCode)
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

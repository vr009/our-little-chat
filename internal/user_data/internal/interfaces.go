package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models2 "our-little-chatik/internal/models"
)

type UserdataRepo interface {
	GetAllUsers() ([]models2.UserData, models2.StatusCode)
	CreateUser(userData models2.UserData) (models2.UserData, models2.StatusCode)
	GetUser(userData models2.UserData) (models2.UserData, models2.StatusCode)
	GetUserForItsName(userData models2.UserData) (models2.UserData, models2.StatusCode)
	DeleteUser(userData models2.UserData) models2.StatusCode
	UpdateUser(personNew models2.UserData) (models2.UserData, models2.StatusCode)
	FindUser(name string) ([]models2.UserData, models2.StatusCode)
}

type UserdataUseCase interface {
	GetAllUsers() ([]models2.UserData, models2.StatusCode)
	CreateUser(userData models2.UserData) (models2.UserData, models2.StatusCode)
	GetUser(userData models2.UserData) (models2.UserData, models2.StatusCode)
	DeleteUser(userData models2.UserData) models2.StatusCode
	UpdateUser(userData models2.UserData) (models2.UserData, models2.StatusCode)
	CheckUser(userData models2.UserData) (models2.UserData, models2.StatusCode)
	FindUser(name string) ([]models2.UserData, models2.StatusCode)
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

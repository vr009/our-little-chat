package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"our-little-chatik/internal/user_data/internal/models"
)

//todo поправить имена в sql-запросах

const (
	INSERTQUERY = "INSERT INTO users(UserID, Nickname, LastAuth, Registered, Avatar, ContactList) " +
		"VALUES($1, $2, $3, $4, $5, $6) RETURNING UserID;"
	DELETEQUERY = "DELETE FROM users WHERE UserID=$1;"
	UPDATEQUERY = "UPDATE users SET Nickname=$1, LastAuth=$2, Registered=$3, Avatar=$4, ContactList=$5 WHERE UserID=$6;"
	GETQUERY    = "SELECT UserID, Nickname, LastAuth, Registered, Avatar, ContactList  FROM users WHERE UserID=$1;"
	LISTQUERY   = "SELECT * FROM users;"
)

type PersonRepo struct {
	conn *pgxpool.Pool
}

func NewPersonRepo(conn *pgxpool.Pool) *PersonRepo {
	return &PersonRepo{
		conn: conn,
	}
}

func (pr *PersonRepo) CreateUser(person models.UserData) (models.UserData, models.StatusCode) {

	_, err := pr.conn.Query(context.Background(),
		INSERTQUERY,
		person.UserID,
		person.Nickname,
		person.LastAuth,
		person.Registered,
		person.Avatar,
		person.ContactList,
	)

	//todo сделать проверку повторного добавления

	if err != nil {
		fmt.Println(err)
		return models.UserData{}, models.InternalError
	}

	return person, models.OK
}

func (pr *PersonRepo) DeleteUser(person models.UserData) models.StatusCode {
	_, err := pr.conn.Exec(context.Background(), DELETEQUERY, person.UserID)
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

	_, err := pr.conn.Exec(context.Background(), UPDATEQUERY, personNew.Nickname, personNew.LastAuth, personNew.Registered, personNew.Avatar, personNew.ContactList, personNew.UserID)
	if err != nil {
		return personOld, models.BadRequest
	}
	personNew = personOld
	return personNew, models.OK
}

func (pr *PersonRepo) GetUser(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.conn.QueryRow(context.Background(), GETQUERY, person.UserID)
	err := rows.Scan(&person.UserID, &person.Nickname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
	if err != nil {
		fmt.Println("ОШИБКА!!!! ", err.Error())
		return models.UserData{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetAllUsers() ([]models.UserData, models.StatusCode) {
	rows, err := pr.conn.Query(context.Background(), LISTQUERY)
	if err != nil && err != sql.ErrNoRows {
		return nil, models.InternalError
	}
	list := make([]models.UserData, 0)
	for rows.Next() {
		person := models.UserData{}
		rows.Scan(&person.UserID, &person.Nickname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
		list = append(list, person)
	}
	return list, models.OK
}

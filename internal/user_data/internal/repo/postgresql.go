package repo

import (
	"context"
	"database/sql"
	"fmt"

	"our-little-chatik/internal/user_data/internal/models"

	"github.com/golang/glog"
	"github.com/jackc/pgx/v4/pgxpool"
)

//todo поправить имена в sql-запросах

const (
	InsertQuery = "INSERT INTO users(user_id, nickname, name, surname, password, last_auth, registered, avatar, contact_list) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);"
	DeleteQuery  = "DELETE FROM users WHERE user_id=$1;"
	UpdateQuery  = "UPDATE users SET nickname=$1, name=$2, surname=$3, last_auth=$4, registered=$5, avatar=$6, contact_list=$7 WHERE user_id=$8;"
	GetQuery     = "SELECT user_id, nickname, name, surname, last_auth, registered, avatar, contact_list  FROM users WHERE user_id=$1;"
	GetNameQuery = "SELECT user_id, nickname, name, surname, password, last_auth, registered, avatar, contact_list  FROM users WHERE nickname=$1;"
	ListQuery    = "SELECT * FROM users;"
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
	glog.Infof("Creating user %s", person)
	_, err := pr.conn.Exec(context.Background(),
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
		glog.Errorf(err.Error())
		return models.UserData{}, models.BadRequest
	}
	person.Password = ""
	return person, models.OK
}

func (pr *PersonRepo) DeleteUser(person models.UserData) models.StatusCode {
	_, err := pr.conn.Exec(context.Background(), DeleteQuery, person.UserID)
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

	_, err := pr.conn.Exec(context.Background(), UpdateQuery, personNew.Nickname, &personNew.Name, &personNew.Surname, personNew.LastAuth, personNew.Registered, personNew.Avatar, personNew.ContactList, personNew.UserID)
	if err != nil {
		return personOld, models.BadRequest
	}
	personNew = personOld
	return personNew, models.OK
}

func (pr *PersonRepo) GetUser(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.conn.QueryRow(context.Background(), GetQuery, person.UserID)
	err := rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
	if err != nil {
		fmt.Println(err)
		return models.UserData{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetUserForItsName(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.conn.QueryRow(context.Background(), GetNameQuery, person.Nickname)
	err := rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.Password, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
	if err != nil {
		fmt.Println(err)
		return models.UserData{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetAllUsers() ([]models.UserData, models.StatusCode) {
	rows, err := pr.conn.Query(context.Background(), ListQuery)
	if err != nil && err != sql.ErrNoRows {
		return nil, models.InternalError
	}
	list := make([]models.UserData, 0)
	for rows.Next() {
		person := models.UserData{}
		rows.Scan(&person.UserID, &person.Nickname, &person.Name, &person.Surname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
		list = append(list, person)
	}
	return list, models.OK
}
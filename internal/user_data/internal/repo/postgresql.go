package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"user_data/internal/models"
)

const (
	INSERTQUERY = "INSERT INTO persons(UserID, Nickname, LastAuth, Registered, Avatar, ContactList) VALUES($1, $2, $3, $4, $5, $6) RETURNING person_id;"
	DELETEQUERY = "DELETE FROM persons WHERE person_id=$1;"
	UPDATEQUERY = "UPDATE persons SET Nickname=$1, LastAuth=$2, Registered=$3, Avatar=$4, ContactList=$5 WHERE person_id=$6;"
	GETQUERY    = "SELECT name, age, work, address FROM persons WHERE person_id=$1;"
	LISTQUERY   = "SELECT * FROM persons;"
)

type PersonRepo struct {
	conn *pgxpool.Pool
}

func NewPersonRepo(conn *pgxpool.Pool) *PersonRepo {
	return &PersonRepo{
		conn: conn,
	}
}

func (pr *PersonRepo) CreatePerson(person models.UserData) (models.UserData, models.StatusCode) {
	NewPerson := models.UserData{UserID: person.UserID, Nickname: person.Nickname, LastAuth: person.LastAuth, Registered: person.Registered, Avatar: person.Avatar, ContactList: person.ContactList}
	row := pr.conn.QueryRow(context.Background(), INSERTQUERY, person.UserID, person.Nickname, person.LastAuth, person.Registered, person.Avatar, person.ContactList)
	err := row.Scan(&NewPerson.UserID)
	if err != nil {
		fmt.Println(err)
		return models.UserData{}, models.BadRequest
	}
	return NewPerson, models.OK
}

func (pr *PersonRepo) DeletePerson(person models.UserData) models.StatusCode {
	_, err := pr.conn.Exec(context.Background(), DELETEQUERY, person.UserID)
	if err != nil {
		return models.InternalError
	}
	return models.OK
}

func (pr *PersonRepo) UpdatePerson(personNew *models.UserData) models.StatusCode {
	personOld, status := pr.GetPerson(*personNew)
	if status != models.OK {
		return models.NotFound
	}

	if personOld.Nickname != personNew.Nickname && personNew.Nickname != "" {
		personOld.Nickname = personNew.Nickname
	}

	if personOld.LastAuth != personNew.LastAuth {
		personOld.LastAuth = personNew.LastAuth
	}

	_, err := pr.conn.Exec(context.Background(), UPDATEQUERY, personOld.Nickname, personOld.LastAuth)
	if err != nil {
		return models.BadRequest
	}
	*personNew = personOld
	return models.OK
}

func (pr *PersonRepo) GetPerson(person models.UserData) (models.UserData, models.StatusCode) {
	rows := pr.conn.QueryRow(context.Background(), GETQUERY, person.UserID)
	err := rows.Scan(&person.Nickname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
	if err != nil {
		return models.UserData{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetPersonsList() ([]models.UserData, models.StatusCode) {
	rows, err := pr.conn.Query(context.Background(), LISTQUERY)
	if err != nil && err != sql.ErrNoRows {
		return nil, models.InternalError
	}
	list := make([]models.UserData, 0)
	for rows.Next() {
		person := models.UserData{}
		rows.Scan(&person.Nickname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
		list = append(list, person)
	}
	return list, models.OK
}

package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"user_data/internal/models"
)

const (
	INSERTQUERY = "INSERT INTO persons(name, age, work, address) VALUES($1, $2, $3, $4) RETURNING person_id;"
	DELETEQUERY = "DELETE FROM persons WHERE person_id=$1;"
	UPDATEQUERY = "UPDATE persons SET name=$1, age=$2, work=$3, address=$4 WHERE person_id=$5;"
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

func (pr *PersonRepo) CreatePerson(person models.Person) (models.Person, models.StatusCode) {
	NewPerson := models.Person{UserId: person.UserId, Nickname: person.Nickname, LastAuth: person.LastAuth, Registered: person.Registered, Avatar: person.Avatar, ContactList: person.ContactList}
	row := pr.conn.QueryRow(context.Background(), INSERTQUERY, person.UserId, person.Nickname, person.LastAuth, person.Registered, person.Avatar, person.ContactList)
	err := row.Scan(&NewPerson.UserId)
	if err != nil {
		fmt.Println(err)
		return models.Person{}, models.BadRequest
	}
	return NewPerson, models.OK
}

func (pr *PersonRepo) DeletePerson(person models.Person) models.StatusCode {
	_, err := pr.conn.Exec(context.Background(), DELETEQUERY, person.UserId)
	if err != nil {
		return models.InternalError
	}
	return models.OK
}

func (pr *PersonRepo) UpdatePerson(personNew *models.Person) models.StatusCode {
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

func (pr *PersonRepo) GetPerson(person models.Person) (models.Person, models.StatusCode) {
	rows := pr.conn.QueryRow(context.Background(), GETQUERY, person.UserId)
	err := rows.Scan(&person.Nickname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
	if err != nil {
		return models.Person{}, models.NotFound
	}
	return person, models.OK
}

func (pr *PersonRepo) GetPersonsList() ([]models.Person, models.StatusCode) {
	rows, err := pr.conn.Query(context.Background(), LISTQUERY)
	if err != nil && err != sql.ErrNoRows {
		return nil, models.InternalError
	}
	list := make([]models.Person, 0)
	for rows.Next() {
		person := models.Person{}
		rows.Scan(&person.Nickname, &person.LastAuth, &person.Registered, &person.Avatar, &person.ContactList)
		list = append(list, person)
	}
	return list, models.OK
}

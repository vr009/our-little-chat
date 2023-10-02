package repo

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"our-little-chatik/internal/models"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestPersonRepo_CreateUser(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		person models.User
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUserID := uuid.New()

	testPassword := "test"

	testRegisterTime := time.Now()

	testPerson := models.User{
		ID:         testUserID,
		Name:       "test",
		Nickname:   "test",
		Surname:    "test",
		Registered: testRegisterTime,
		Avatar:     "avatar.png",
	}
	err = testPerson.Password.Set(testPassword)
	if err != nil {
		t.Fatal(err.Error())
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.User
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(InsertQuery)).
					WithArgs(testPerson.ID, testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Password.Hash,
						testPerson.Avatar).
					WillReturnRows(sqlmock.NewRows([]string{"registered"}).AddRow(testPerson.Registered))
			},
			fields: fields{
				pool: db,
			},
			args: args{
				person: testPerson,
			},
			want: models.User{
				ID:         testUserID,
				Name:       "test",
				Nickname:   "test",
				Surname:    "test",
				Password:   testPerson.Password,
				Registered: testRegisterTime,
				Avatar:     "avatar.png",
			},
			want1: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PersonRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, got1 := pr.CreateUser(tt.args.person)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CreateUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPersonRepo_DeleteUser(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		person models.User
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()

	testPerson := models.User{
		ID:         testUserID,
		Name:       "test",
		Nickname:   "test",
		Surname:    "test",
		Registered: time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}

	tests := []struct {
		name   string
		pre    func()
		fields fields
		args   args
		want   models.StatusCode
	}{
		{
			name: "",
			fields: fields{
				pool: db,
			},
			args: args{
				person: testPerson,
			},
			want: models.Deleted,
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(DeleteQuery)).
					WithArgs(testPerson.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PersonRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			if got := pr.DeleteUser(tt.args.person); got != tt.want {
				t.Errorf("DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPersonRepo_FindUser(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		name string
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUserID := uuid.New()

	testPerson := models.User{
		ID:       testUserID,
		Name:     "test",
		Nickname: "test",
		Surname:  "test",
		Avatar:   "avatar.png",
	}

	columns := []string{
		"user_id", "nickname", "name", "surname", "avatar",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   []models.User
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(FindUsersQuery)).
					WithArgs(testPerson.Nickname).WillReturnRows(sqlmock.NewRows(columns).
					AddRow(testPerson.ID.String(), testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Avatar))
			},
			fields: fields{pool: db},
			args: args{
				name: testPerson.Nickname,
			},
			want:  []models.User{testPerson},
			want1: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PersonRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, got1 := pr.FindUser(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPersonRepo_GetAllUsers(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUserID := uuid.New()

	testPerson := models.User{
		ID:       testUserID,
		Nickname: "test",
		Avatar:   "avatar.png",
	}

	columns := []string{
		"user_id", "nickname", "avatar",
	}

	tests := []struct {
		name   string
		pre    func()
		fields fields
		want   []models.User
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(ListQuery)).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(testPerson.ID.String(), testPerson.Nickname, testPerson.Avatar))
			},
			fields: fields{pool: db},
			want:   []models.User{testPerson},
			want1:  models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PersonRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, got1 := pr.GetAllUsers()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllUsers() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAllUsers() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPersonRepo_GetUser(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		person models.User
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()
	password := "test"

	testPerson := models.User{
		ID:         testUserID,
		Name:       "test",
		Nickname:   "test",
		Surname:    "test",
		Registered: time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}

	testPerson.Password.Set(password)

	columns := []string{
		"user_id", "nickname", "name", "surname", "password", "registered",
		"avatar",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.User
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetQuery)).
					WithArgs(testPerson.ID).WillReturnRows(sqlmock.NewRows(columns).
					AddRow(testPerson.ID.String(), testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Password.Hash,
						testPerson.Registered, testPerson.Avatar))
			},
			fields: fields{pool: db},
			args: args{
				person: testPerson,
			},
			want:  testPerson,
			want1: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PersonRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, got1 := pr.GetUser(tt.args.person)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPersonRepo_GetUserForItsName(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		person models.User
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()
	testPassword := "test"

	testPerson := models.User{
		ID:         testUserID,
		Name:       "test",
		Nickname:   "test",
		Surname:    "test",
		Registered: time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}
	testPerson.Password.Set(testPassword)

	columns := []string{
		"user_id", "nickname", "name", "surname", "password", "registered",
		"avatar",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.User
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetNameQuery)).
					WithArgs(testPerson.Name).WillReturnRows(sqlmock.NewRows(columns).
					AddRow(testPerson.ID.String(), testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Password.Hash, testPerson.Registered, testPerson.Avatar))
			},
			fields: fields{pool: db},
			args: args{
				person: testPerson,
			},
			want:  testPerson,
			want1: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PersonRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, got1 := pr.GetUserForItsName(tt.args.person)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserForItsName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetUserForItsName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPersonRepo_UpdateUser(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		personNew models.User
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testPassword := "test"
	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()

	testPerson := models.User{
		ID:         testUserID,
		Name:       "test",
		Nickname:   "test",
		Surname:    "test",
		Registered: time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}

	testPerson.Password.Set(testPassword)

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.User
		want1  models.StatusCode
	}{
		{
			name: "",
			fields: fields{
				pool: db,
			},
			args: args{
				personNew: testPerson,
			},
			want: testPerson,
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(UpdateQuery)).
					WithArgs(testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Avatar, testPerson.Password.Hash, testPerson.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &PersonRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, got1 := pr.UpdateUser(tt.args.personNew)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UpdateUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

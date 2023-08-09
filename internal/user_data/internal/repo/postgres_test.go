package repo

import (
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestPersonRepo_CreateUser(t *testing.T) {
	type fields struct {
		pool internal.DB
	}
	type args struct {
		person models.UserData
	}
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()

	testPerson := models.UserData{
		User: models2.User{
			UserID:   testUserID,
			Name:     "test",
			Nickname: "test",
			Surname:  "test",
			Password: "test",
		},
		Registered: time.Unix(testTimestamp, 0),
		LastAuth:   time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.UserData
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(InsertQuery)).
					WithArgs(testPerson.UserID, testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Password, testPerson.LastAuth,
						testPerson.Registered, testPerson.Avatar).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			fields: fields{
				pool: mock,
			},
			args: args{
				person: testPerson,
			},
			want: models.UserData{
				User: models2.User{
					UserID:   testUserID,
					Name:     "test",
					Nickname: "test",
					Surname:  "test",
				},
				Registered: time.Unix(testTimestamp, 0),
				LastAuth:   time.Unix(testTimestamp, 0),
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
		pool internal.DB
	}
	type args struct {
		person models.UserData
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()

	testPerson := models.UserData{
		User: models2.User{
			UserID:   testUserID,
			Name:     "test",
			Nickname: "test",
			Surname:  "test",
			Password: "test",
		},
		Registered: time.Unix(testTimestamp, 0),
		LastAuth:   time.Unix(testTimestamp, 0),
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
				pool: mock,
			},
			args: args{
				person: testPerson,
			},
			want: models.Deleted,
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(DeleteQuery)).
					WithArgs(testPerson.UserID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
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
		pool internal.DB
	}
	type args struct {
		name string
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()

	testPerson := models.UserData{
		User: models2.User{
			UserID:   testUserID,
			Name:     "test",
			Nickname: "test",
			Surname:  "test",
		},
		Avatar: "avatar.png",
	}

	columns := []string{
		"user_id", "nickname", "name", "surname", "avatar",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   []models.UserData
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(FindUsersQuery)).
					WithArgs(testPerson.Nickname).WillReturnRows(pgxmock.NewRows(columns).
					AddRow(testPerson.UserID.String(), testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Avatar))
			},
			fields: fields{pool: mock},
			args: args{
				name: testPerson.Nickname,
			},
			want:  []models.UserData{testPerson},
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
		pool internal.DB
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()

	testPerson := models.UserData{
		User: models2.User{
			UserID:   testUserID,
			Nickname: "test",
		},
		Avatar: "avatar.png",
	}

	columns := []string{
		"user_id", "nickname", "avatar",
	}

	tests := []struct {
		name   string
		pre    func()
		fields fields
		want   []models.UserData
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(ListQuery)).
					WillReturnRows(pgxmock.NewRows(columns).
						AddRow(testPerson.UserID.String(), testPerson.Nickname, testPerson.Avatar))
			},
			fields: fields{pool: mock},
			want:   []models.UserData{testPerson},
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
		pool internal.DB
	}
	type args struct {
		person models.UserData
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()

	testPerson := models.UserData{
		User: models2.User{
			UserID:   testUserID,
			Name:     "test",
			Nickname: "test",
			Surname:  "test",
			Password: "test",
		},
		Registered: time.Unix(testTimestamp, 0),
		LastAuth:   time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}

	columns := []string{
		"user_id", "nickname", "name", "surname", "last_auth", "registered",
		"avatar",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.UserData
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetQuery)).
					WithArgs(testPerson.UserID).WillReturnRows(pgxmock.NewRows(columns).
					AddRow(testPerson.UserID.String(), testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.LastAuth, testPerson.Registered, testPerson.Avatar))
			},
			fields: fields{pool: mock},
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
		pool internal.DB
	}
	type args struct {
		person models.UserData
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()

	testPerson := models.UserData{
		User: models2.User{
			UserID:   testUserID,
			Name:     "test",
			Nickname: "test",
			Surname:  "test",
			Password: "test",
		},
		Registered: time.Unix(testTimestamp, 0),
		LastAuth:   time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}

	columns := []string{
		"user_id", "nickname", "name", "surname", "last_auth", "registered",
		"avatar",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.UserData
		want1  models.StatusCode
	}{
		{
			name: "",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetNameQuery)).
					WithArgs(testPerson.Name).WillReturnRows(pgxmock.NewRows(columns).
					AddRow(testPerson.UserID.String(), testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.LastAuth, testPerson.Registered, testPerson.Avatar))
			},
			fields: fields{pool: mock},
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
		pool internal.DB
	}
	type args struct {
		personNew models.UserData
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()
	testTimestamp := time.Now().Unix()

	testPerson := models.UserData{
		User: models2.User{
			UserID:   testUserID,
			Name:     "test",
			Nickname: "test",
			Surname:  "test",
			Password: "test",
		},
		Registered: time.Unix(testTimestamp, 0),
		LastAuth:   time.Unix(testTimestamp, 0),
		Avatar:     "avatar.png",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.UserData
		want1  models.StatusCode
	}{
		{
			name: "",
			fields: fields{
				pool: mock,
			},
			args: args{
				personNew: testPerson,
			},
			want: testPerson,
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(UpdateQuery)).
					WithArgs(testPerson.Nickname, testPerson.Name,
						testPerson.Surname, testPerson.Password, testPerson.LastAuth,
						testPerson.Registered, testPerson.Avatar, testPerson.UserID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
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

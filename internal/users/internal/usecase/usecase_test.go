package usecase

import (
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"our-little-chatik/internal/models"
	mocks "our-little-chatik/internal/users/internal/mocks/users"
	models2 "our-little-chatik/internal/users/internal/models"
	"reflect"
	"testing"
)

func TestUserUsecase_Login(t *testing.T) {
	type fields struct {
		repo *mocks.MockUserRepo
	}
	type args struct {
		request models2.LoginRequest
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testNickname := "test"

	testEmptyUser := models.User{}

	testUser := models.User{
		ID:        uuid.New(),
		Name:      "test",
		Nickname:  testNickname,
		Surname:   "test",
		Avatar:    "test",
		Activated: true,
	}

	testInActivatedUser := models.User{
		ID:        uuid.New(),
		Name:      "test",
		Nickname:  testNickname,
		Surname:   "test",
		Avatar:    "test",
		Activated: false,
	}

	testPassword := "testPswd"
	err := testUser.Password.Set(testPassword)
	if err != nil {
		t.Fatal(err)
	}

	testBadPassword := "testWrongPassword"

	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(f *fields)
		want    models.User
		want1   models.StatusCode
	}{
		{
			name: "successful login",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				models2.LoginRequest{
					Password: &testPassword,
					Nickname: &testNickname,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().
					GetUserForItsNickname(models.User{Nickname: testNickname}).
					Return(testUser, models.OK)
			},
			want:  testUser,
			want1: models.OK,
		},
		{
			name: "user can not be found",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				models2.LoginRequest{
					Password: &testPassword,
					Nickname: &testNickname,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().
					GetUserForItsNickname(models.User{Nickname: testNickname}).
					Return(testEmptyUser, models.NotFound)
			},
			want:  testEmptyUser,
			want1: models.NotFound,
		},
		{
			name: "bad credentials",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				models2.LoginRequest{
					Password: &testBadPassword,
					Nickname: &testNickname,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().
					GetUserForItsNickname(models.User{Nickname: testNickname}).
					Return(testUser, models.OK)
			},
			want:  testEmptyUser,
			want1: models.Unauthorized,
		},
		{
			name: "bad credentials",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				models2.LoginRequest{
					Password: &testBadPassword,
					Nickname: &testNickname,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().
					GetUserForItsNickname(models.User{Nickname: testNickname}).
					Return(testInActivatedUser, models.OK)
			},
			want:  testEmptyUser,
			want1: models.InActivated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &UserUsecase{
				userRepo: tt.fields.repo,
			}
			tt.prepare(&tt.fields)
			got, got1 := uc.Login(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Login() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserUsecase_SignUp(t *testing.T) {
	type fields struct {
		repo *mocks.MockUserRepo
	}
	type args struct {
		request models2.SignUpPersonRequest
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testNickname := "test"

	testEmptyUser := models.User{}

	testUser := models.User{
		ID:        uuid.New(),
		Name:      "test",
		Nickname:  testNickname,
		Surname:   "test",
		Avatar:    "test",
		Activated: true,
	}

	testPassword := "testPswd"
	err := testUser.Password.Set(testPassword)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(f *fields)
		want    models.User
		want1   models.StatusCode
	}{
		{
			name: "successful sign up",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				models2.SignUpPersonRequest{
					Nickname: &testUser.Nickname,
					Name:     &testUser.Name,
					Surname:  &testUser.Surname,
					Avatar:   &testUser.Avatar,
					Password: &testPassword,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().CreateUser(gomock.Cond(func(x any) bool {
					usr, ok := x.(models.User)
					if !ok {
						return false
					}
					ok, _ = testUser.Password.Matches(*usr.Password.Plaintext)
					return usr.Name == testUser.Name &&
						usr.Nickname == testUser.Nickname &&
						usr.Surname == testUser.Surname &&
						usr.Activated == testUser.Activated &&
						usr.Avatar == testUser.Avatar && ok
				})).Return(testUser, models.OK)
			},
			want:  testUser,
			want1: models.OK,
		},
		{
			name: "conflict",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				models2.SignUpPersonRequest{
					Nickname: &testUser.Nickname,
					Name:     &testUser.Name,
					Surname:  &testUser.Surname,
					Avatar:   &testUser.Avatar,
					Password: &testPassword,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().CreateUser(gomock.Cond(func(x any) bool {
					usr, ok := x.(models.User)
					if !ok {
						return false
					}
					ok, _ = testUser.Password.Matches(*usr.Password.Plaintext)
					return usr.Name == testUser.Name &&
						usr.Nickname == testUser.Nickname &&
						usr.Surname == testUser.Surname &&
						usr.Activated == testUser.Activated &&
						usr.Avatar == testUser.Avatar && ok
				})).Return(testEmptyUser, models.Conflict)
			},
			want:  testEmptyUser,
			want1: models.Conflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := &UserUsecase{
				userRepo: tt.fields.repo,
			}
			tt.prepare(&tt.fields)
			got, got1 := uc.SignUp(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignUp() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SignUp() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	type fields struct {
		repo *mocks.MockUserRepo
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testNickname := "test"
	testModifyNickname := "testModified"

	//testEmptyUser := models.User{}

	testUser := models.User{
		ID:        uuid.New(),
		Name:      "test",
		Nickname:  testNickname,
		Surname:   "test",
		Avatar:    "test",
		Activated: true,
	}

	testInActivatedUser := models.User{
		ID:        uuid.New(),
		Name:      "test",
		Nickname:  testNickname,
		Surname:   "test",
		Avatar:    "test",
		Activated: false,
	}

	testModifiedNicknameUser := models.User{
		ID:        testUser.ID,
		Nickname:  testModifyNickname,
		Name:      testUser.Name,
		Surname:   testUser.Surname,
		Avatar:    testUser.Avatar,
		Activated: testUser.Activated,
	}

	testPassword := "testPswd"
	err := testUser.Password.Set(testPassword)
	if err != nil {
		t.Fatal(err)
	}

	testModifiedNicknameUser.Password.Hash = testUser.Password.Hash

	type args struct {
		userToUpdate models.User
		request      models2.UpdateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		prepare func(f *fields)
		want    models.User
		want1   models.StatusCode
	}{
		{
			name: "successful update",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				userToUpdate: models.User{ID: testUser.ID},
				request: models2.UpdateUserRequest{
					Nickname: &testModifyNickname,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().GetUserForItsID(models.User{ID: testUser.ID}).
					Return(testUser, models.OK)
				f.repo.EXPECT().UpdateUser(testModifiedNicknameUser).
					Return(testModifiedNicknameUser, models.OK)
			},
			want:  testModifiedNicknameUser,
			want1: models.OK,
		},
		{
			name: "failure by update inactive user",
			fields: fields{
				repo: mocks.NewMockUserRepo(ctrl),
			},
			args: args{
				userToUpdate: models.User{ID: testInActivatedUser.ID},
				request: models2.UpdateUserRequest{
					Nickname: &testModifyNickname,
				},
			},
			prepare: func(f *fields) {
				f.repo.EXPECT().GetUserForItsID(models.User{ID: testInActivatedUser.ID}).
					Return(testInActivatedUser, models.OK)
			},
			want:  models.User{},
			want1: models.InActivated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc := &UserUsecase{
				userRepo: tt.fields.repo,
			}
			tt.prepare(&tt.fields)
			got, got1 := uc.UpdateUser(tt.args.userToUpdate, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UpdateUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

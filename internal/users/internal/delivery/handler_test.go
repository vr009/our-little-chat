package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"our-little-chatik/internal/models"
	mocks "our-little-chatik/internal/users/internal/mocks/users"
	models2 "our-little-chatik/internal/users/internal/models"
	"strings"
	"testing"
)

func TestUserEchoHandler_UpdateUser(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}
	type args struct {
		c echo.Context
	}

	testID := uuid.New()
	testNickname := "testNick"
	testName := "testName"
	testSurname := "testSurname"
	testPassword := "testpassword"
	testAvatar := "avatar.png"

	testNewNickname := "newtestNick"
	testNewName := "newtestName"
	testNewSurname := "newtestSurname"
	testNewPassword := "newtestpassword"
	testNewAvatar := "new_avatar.png"

	testEmptyUser := models.User{}

	testUser := models.User{
		ID:        testID,
		Name:      testName,
		Nickname:  testNickname,
		Surname:   testSurname,
		Avatar:    testAvatar,
		Activated: true,
	}
	testUser.Password.Set(testPassword)

	testNewUser := models.User{
		ID:        testID,
		Name:      testNewName,
		Nickname:  testNewNickname,
		Surname:   testNewSurname,
		Avatar:    testNewAvatar,
		Activated: true,
	}
	testNewUser.Password.Set(testNewPassword)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                 string
		fields               fields
		args                 args
		prepareEchoCtx       func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder)
		prepareUpdateRequest func() models2.UpdateUserRequest
		prepare              func(f *fields, input models2.UpdateUserRequest)
		checkResponse        func(recorder *httptest.ResponseRecorder) error
		wantErr              bool
	}{
		{
			name: "successful update - all fields",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareUpdateRequest: func() models2.UpdateUserRequest {
				testInput := models2.UpdateUserRequest{
					Nickname:    &testNewNickname,
					Name:        &testNewName,
					Surname:     &testNewSurname,
					OldPassword: &testPassword,
					NewPassword: &testNewPassword,
					Avatar:      &testNewAvatar,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", testID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.UpdateUserRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().UpdateUser(models.User{ID: testID}, input).Return(testNewUser, models.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "successful update - only password",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareUpdateRequest: func() models2.UpdateUserRequest {
				testInput := models2.UpdateUserRequest{
					OldPassword: &testPassword,
					NewPassword: &testNewPassword,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", testID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.UpdateUserRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().UpdateUser(models.User{ID: testID}, input).Return(testNewUser, models.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "not successful update - old password is not provided",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareUpdateRequest: func() models2.UpdateUserRequest {
				testInput := models2.UpdateUserRequest{
					NewPassword: &testNewPassword,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", testID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.UpdateUserRequest) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnprocessableEntity {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "not successful update - new password is not provided",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareUpdateRequest: func() models2.UpdateUserRequest {
				testInput := models2.UpdateUserRequest{
					OldPassword: &testPassword,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", testID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.UpdateUserRequest) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnprocessableEntity {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "not successful update - bad body",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareUpdateRequest: func() models2.UpdateUserRequest {
				testInput := models2.UpdateUserRequest{
					OldPassword: &testPassword,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&[]string{"dd"})
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", testID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.UpdateUserRequest) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusBadRequest {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "not successful update - user not found",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareUpdateRequest: func() models2.UpdateUserRequest {
				testInput := models2.UpdateUserRequest{
					Nickname:    &testNewNickname,
					Name:        &testNewName,
					Surname:     &testNewSurname,
					OldPassword: &testPassword,
					NewPassword: &testNewPassword,
					Avatar:      &testNewAvatar,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", testID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.UpdateUserRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().UpdateUser(models.User{ID: testID}, input).Return(testEmptyUser, models.NotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusNotFound {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "not successful update - failed to update",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareUpdateRequest: func() models2.UpdateUserRequest {
				testInput := models2.UpdateUserRequest{
					Nickname:    &testNewNickname,
					Name:        &testNewName,
					Surname:     &testNewSurname,
					OldPassword: &testPassword,
					NewPassword: &testNewPassword,
					Avatar:      &testNewAvatar,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.UpdateUserRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", testID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.UpdateUserRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().UpdateUser(models.User{ID: testID}, input).Return(testEmptyUser, models.InternalError)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusInternalServerError {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udh := &UserEchoHandler{
				useCase: tt.fields.useCase,
			}
			input := tt.prepareUpdateRequest()
			tt.prepare(&tt.fields, input)
			c, rec := tt.prepareEchoCtx(input)
			if err := udh.UpdateUser(c); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUserEchoHandler_GetMe(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}
	type args struct {
		c echo.Context
	}

	testID := uuid.New()
	testUser := models.User{
		ID:       testID,
		Name:     "test",
		Surname:  "test",
		Nickname: "test",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	testEchoCtx := e.NewContext(req, rec)
	testEchoCtx.Set("user_id", testID)

	testID2 := uuid.New()
	rec2 := httptest.NewRecorder()
	testEchoCtx2 := e.NewContext(req, rec2)
	testEchoCtx2.Set("user_id", testID2)

	testID3 := uuid.New()
	rec3 := httptest.NewRecorder()
	testEchoCtx3 := e.NewContext(req, rec3)
	testEchoCtx3.Set("user_id", testID3)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		fields        fields
		prepare       func(f *fields)
		args          args
		checkResponse func(recorder *httptest.ResponseRecorder) error
		wantErr       bool
		rec           *httptest.ResponseRecorder
	}{
		{
			name: "successful get",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{
				c: testEchoCtx,
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().GetUser(models2.GetUserRequest{testID}).Return(testUser, models.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
			rec:     rec,
		},
		{
			name: "not successful get - user not found",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{
				c: testEchoCtx2,
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().GetUser(models2.GetUserRequest{testID2}).Return(models.User{}, models.NotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusNotFound {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
			rec:     rec2,
		},
		{
			name: "not successful get - internal issue",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{
				c: testEchoCtx3,
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().GetUser(models2.GetUserRequest{testID3}).Return(models.User{}, models.InternalError)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusInternalServerError {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
			rec:     rec3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udh := &UserEchoHandler{
				useCase: tt.fields.useCase,
			}
			tt.prepare(&tt.fields)
			if err := udh.GetMe(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetMe() error = %v, wantErr %v", err, tt.wantErr)
			}
			err := tt.checkResponse(tt.rec)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUserEchoHandler_GetUserForID(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c echo.Context
	}

	testID := uuid.New()
	testUser := models.User{
		ID:       testID,
		Name:     "test",
		Surname:  "test",
		Nickname: "test",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	testEchoCtx := e.NewContext(req, rec)
	testEchoCtx.Set("user_id", testID)
	testEchoCtx.SetParamNames("id")
	testEchoCtx.SetParamValues(testID.String())

	rec2 := httptest.NewRecorder()
	testEchoCtx2 := e.NewContext(req, rec2)
	testEchoCtx2.Set("user_id", testID)
	testEchoCtx2.SetParamNames("id")
	testEchoCtx2.SetParamValues("badID")

	rec3 := httptest.NewRecorder()
	testEchoCtx3 := e.NewContext(req, rec3)
	testEchoCtx3.Set("user_id", testID)
	testEchoCtx3.SetParamNames("id")
	testEchoCtx3.SetParamValues(testID.String())

	rec4 := httptest.NewRecorder()
	testEchoCtx4 := e.NewContext(req, rec4)
	testEchoCtx4.Set("user_id", testID)
	testEchoCtx4.SetParamNames("id")
	testEchoCtx4.SetParamValues(testID.String())

	tests := []struct {
		name          string
		fields        fields
		args          args
		checkResponse func(recorder *httptest.ResponseRecorder) error
		prepare       func(f *fields)
		rec           *httptest.ResponseRecorder
		wantErr       bool
	}{
		{
			name: "successful - user found",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{c: testEchoCtx},
			prepare: func(f *fields) {
				f.useCase.EXPECT().GetUser(models2.GetUserRequest{testID}).Return(testUser, models.OK)
			},
			rec:     rec,
			wantErr: false,
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
		},
		{
			name: "failure - bad id param",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{c: testEchoCtx2},
			prepare: func(f *fields) {
			},
			rec:     rec2,
			wantErr: false,
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusBadRequest {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
		},
		{
			name: "failure - user not found",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{c: testEchoCtx3},
			prepare: func(f *fields) {
				f.useCase.EXPECT().GetUser(models2.GetUserRequest{testID}).Return(models.User{}, models.NotFound)
			},
			rec:     rec3,
			wantErr: false,
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusNotFound {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
		},
		{
			name: "failure - internal issue by searching",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{c: testEchoCtx4},
			prepare: func(f *fields) {
				f.useCase.EXPECT().GetUser(models2.GetUserRequest{testID}).Return(models.User{}, models.InternalError)
			},
			rec:     rec4,
			wantErr: false,
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusInternalServerError {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udh := &UserEchoHandler{
				useCase: tt.fields.useCase,
			}
			tt.prepare(&tt.fields)
			if err := udh.GetUserForID(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetUserForID() error = %v, wantErr %v", err, tt.wantErr)
			}
			err := tt.checkResponse(tt.rec)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUserEchoHandler_SearchUsers(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}
	type args struct {
		c echo.Context
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testID := uuid.New()
	testUser := models.User{
		ID:       testID,
		Name:     "test",
		Surname:  "test",
		Nickname: "test",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	testEchoCtx := e.NewContext(req, rec)
	testEchoCtx.Set("user_id", testID)
	testEchoCtx.QueryParams().Set("nickname", testUser.Nickname)

	rec2 := httptest.NewRecorder()
	testEchoCtx2 := e.NewContext(req, rec2)
	testEchoCtx2.Set("user_id", testID)
	testEchoCtx2.QueryParams().Set("nickname", testUser.Nickname)

	rec3 := httptest.NewRecorder()
	testEchoCtx3 := e.NewContext(req, rec3)
	testEchoCtx3.Set("user_id", testID)
	testEchoCtx3.QueryParams().Set("nickname", testUser.Nickname)

	rec4 := httptest.NewRecorder()
	testEchoCtx4 := e.NewContext(req, rec4)
	testEchoCtx4.Set("user_id", testID)

	tests := []struct {
		name          string
		fields        fields
		checkResponse func(recorder *httptest.ResponseRecorder) error
		prepare       func(f *fields)
		rec           *httptest.ResponseRecorder
		args          args
		wantErr       bool
	}{
		{
			name:   "successful search",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().FindUsers(testUser.Nickname).
					Return([]models.User{testUser}, models.OK)
			},
			rec:     rec,
			args:    args{c: testEchoCtx},
			wantErr: false,
		},
		{
			name:   "failure by search - not found",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusNotFound {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().FindUsers(testUser.Nickname).
					Return([]models.User{}, models.NotFound)
			},
			rec:     rec2,
			args:    args{c: testEchoCtx2},
			wantErr: false,
		},
		{
			name:   "failure by search - internal issue",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusInternalServerError {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().FindUsers(testUser.Nickname).
					Return([]models.User{}, models.InternalError)
			},
			rec:     rec3,
			args:    args{c: testEchoCtx3},
			wantErr: false,
		},
		{
			name:   "failure by search - validation error",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnprocessableEntity {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
			},
			rec:     rec4,
			args:    args{c: testEchoCtx4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udh := &UserEchoHandler{
				useCase: tt.fields.useCase,
			}
			tt.prepare(&tt.fields)
			if err := udh.SearchUsers(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SearchUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
			err := tt.checkResponse(tt.rec)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUserEchoHandler_DeactivateUser(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	testEchoCtx := e.NewContext(req, rec)
	testEchoCtx.Set("user_id", testID)

	rec2 := httptest.NewRecorder()
	testEchoCtx2 := e.NewContext(req, rec2)
	testEchoCtx2.Set("user_id", testID)

	rec3 := httptest.NewRecorder()
	testEchoCtx3 := e.NewContext(req, rec3)
	testEchoCtx3.Set("user_id", testID)

	rec4 := httptest.NewRecorder()
	testEchoCtx4 := e.NewContext(req, rec4)

	type args struct {
		c echo.Context
	}
	tests := []struct {
		name          string
		fields        fields
		checkResponse func(recorder *httptest.ResponseRecorder) error
		prepare       func(f *fields)
		rec           *httptest.ResponseRecorder
		args          args
		wantErr       bool
	}{
		{
			name:   "successful search",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().DeactivateUser(models.User{ID: testID}).
					Return(models.OK)
			},
			rec:     rec,
			args:    args{c: testEchoCtx},
			wantErr: false,
		},
		{
			name:   "failure by search - not found",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusNotFound {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().DeactivateUser(models.User{ID: testID}).
					Return(models.NotFound)
			},
			rec:     rec2,
			args:    args{c: testEchoCtx2},
			wantErr: false,
		},
		{
			name:   "failure by search - internal issue",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusInternalServerError {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
				f.useCase.EXPECT().DeactivateUser(models.User{ID: testID}).
					Return(models.InternalError)
			},
			rec:     rec3,
			args:    args{c: testEchoCtx3},
			wantErr: false,
		},
		{
			name:   "failure by search - validation error",
			fields: fields{useCase: mocks.NewMockUserUsecase(ctrl)},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusBadRequest {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields) {
			},
			rec:     rec4,
			args:    args{c: testEchoCtx4},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udh := &UserEchoHandler{
				useCase: tt.fields.useCase,
			}
			tt.prepare(&tt.fields)
			if err := udh.DeactivateUser(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("DeactivateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			err := tt.checkResponse(tt.rec)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

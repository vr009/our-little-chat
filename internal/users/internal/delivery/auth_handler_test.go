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

func TestAuthEchoHandler_Login(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}
	type args struct {
		c echo.Context
	}

	testShortPswd := "test"
	testOkPswd := "testtesttest"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                string
		fields              fields
		prepareEchoCtx      func(input models2.LoginRequest) (echo.Context, *httptest.ResponseRecorder)
		prepareLoginRequest func() models2.LoginRequest
		prepare             func(f *fields, input models2.LoginRequest)
		checkResponse       func(recorder *httptest.ResponseRecorder) error
		args                args
		wantErr             bool
	}{
		{
			name: "successful login",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.LoginRequest {
				testInput := models2.LoginRequest{
					Password: &testOkPswd,
					Nickname: &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.LoginRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.LoginRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().Login(input).Return(testUser, models.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusSeeOther {
					return fmt.Errorf("wrong status code")
				}
				if len(recorder.Result().Cookies()) == 0 {
					return fmt.Errorf("no cookies provided")
				}
				if recorder.Result().Cookies()[0].Name != "Token" {
					return fmt.Errorf("unknown cookie provided")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "not successful login - short pswd",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.LoginRequest {
				testInput := models2.LoginRequest{
					Password: &testShortPswd,
					Nickname: &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.LoginRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnprocessableEntity {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields, input models2.LoginRequest) {
			},
			wantErr: false,
		},
		{
			name: "not successful login - empty body sent",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.LoginRequest {
				testInput := models2.LoginRequest{}
				return testInput
			},
			prepareEchoCtx: func(input models2.LoginRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnprocessableEntity {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			prepare: func(f *fields, input models2.LoginRequest) {
			},
			wantErr: false,
		},
		{
			name: "not successful login - user not found",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.LoginRequest {
				testInput := models2.LoginRequest{
					Password: &testOkPswd,
					Nickname: &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.LoginRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.LoginRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().Login(input).Return(testEmptyUser, models.NotFound)
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
			name: "not successful login - invalid creds",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.LoginRequest {
				testInput := models2.LoginRequest{
					Password: &testOkPswd,
					Nickname: &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.LoginRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.LoginRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().Login(input).Return(testEmptyUser, models.Conflict)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnauthorized {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "not successful login - bad body",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.LoginRequest {
				testInput := models2.LoginRequest{
					Password: &testOkPswd,
					Nickname: &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.LoginRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&[]string{"ewefwe"})
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.LoginRequest) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusBadRequest {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthEchoHandler{
				useCase: tt.fields.useCase,
			}
			input := tt.prepareLoginRequest()
			tt.prepare(&tt.fields, input)
			c, rec := tt.prepareEchoCtx(input)
			if err := h.Login(c); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAuthEchoHandler_Logout(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type args struct {
		c echo.Context
	}

	t.Setenv("JWT_SIGNED_KEY", "test")

	testID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	testEchoCtx := e.NewContext(req, rec)
	testEchoCtx.Set("user_id", testID)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful logout",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			args: args{
				c: testEchoCtx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthEchoHandler{
				useCase: tt.fields.useCase,
			}
			if err := h.Logout(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthEchoHandler_SignUp(t *testing.T) {
	type fields struct {
		useCase *mocks.MockUserUsecase
	}

	testShortPswd := "test"
	testOkPswd := "testtesttest"
	testNickname := "test"

	testUser := models.User{
		ID:        uuid.New(),
		Name:      "test",
		Nickname:  testNickname,
		Surname:   "test",
		Avatar:    "test",
		Activated: true,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		c echo.Context
	}
	tests := []struct {
		name                string
		fields              fields
		prepareEchoCtx      func(input models2.SignUpPersonRequest) (echo.Context, *httptest.ResponseRecorder)
		prepareLoginRequest func() models2.SignUpPersonRequest
		prepare             func(f *fields, input models2.SignUpPersonRequest)
		checkResponse       func(recorder *httptest.ResponseRecorder) error
		args                args
		wantErr             bool
	}{
		{
			name: "successful sign up",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.SignUpPersonRequest {
				testInput := models2.SignUpPersonRequest{
					Password: &testOkPswd,
					Nickname: &testNickname,
					Name:     &testNickname,
					Surname:  &testNickname,
					Avatar:   &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.SignUpPersonRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.SignUpPersonRequest) {
				t.Setenv("JWT_SIGNED_KEY", "test")
				f.useCase.EXPECT().SignUp(input).Return(testUser, models.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusSeeOther {
					return fmt.Errorf("wrong status code")
				}
				if len(recorder.Result().Cookies()) == 0 {
					return fmt.Errorf("no cookies provided")
				}
				if recorder.Result().Cookies()[0].Name != "Token" {
					return fmt.Errorf("unknown cookie provided")
				}
				return nil
			},
		},
		{
			name: "not successful sign up - no password provided",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.SignUpPersonRequest {
				testInput := models2.SignUpPersonRequest{
					Nickname: &testNickname,
					Name:     &testNickname,
					Surname:  &testNickname,
					Avatar:   &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.SignUpPersonRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.SignUpPersonRequest) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnprocessableEntity {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
		},
		{
			name: "not successful sign up - short password provided",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.SignUpPersonRequest {
				testInput := models2.SignUpPersonRequest{
					Nickname: &testNickname,
					Name:     &testNickname,
					Password: &testShortPswd,
					Surname:  &testNickname,
					Avatar:   &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.SignUpPersonRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.SignUpPersonRequest) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnprocessableEntity {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
		},
		{
			name: "not successful sign up - user exists",
			fields: fields{
				useCase: mocks.NewMockUserUsecase(ctrl),
			},
			prepareLoginRequest: func() models2.SignUpPersonRequest {
				testInput := models2.SignUpPersonRequest{
					Nickname: &testNickname,
					Name:     &testNickname,
					Password: &testOkPswd,
					Surname:  &testNickname,
					Avatar:   &testNickname,
				}
				return testInput
			},
			prepareEchoCtx: func(input models2.SignUpPersonRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models2.SignUpPersonRequest) {
				f.useCase.EXPECT().SignUp(input).Return(testUser, models.Conflict)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusUnauthorized {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthEchoHandler{
				useCase: tt.fields.useCase,
			}
			input := tt.prepareLoginRequest()
			tt.prepare(&tt.fields, input)
			c, rec := tt.prepareEchoCtx(input)
			if err := h.SignUp(c); (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err)
			}
		})
	}
}

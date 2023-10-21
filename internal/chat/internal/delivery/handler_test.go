package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"our-little-chatik/internal/chat/internal/mocks/chat"
	"our-little-chatik/internal/chat/internal/models"
	models2 "our-little-chatik/internal/models"
	"strings"
	"testing"
)

func TestChatEchoHandler_AddUsersToChat(t *testing.T) {
	type fields struct {
		usecase *chat.MockChatUseCase
	}
	type args struct {
		c echo.Context
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatID := uuid.New()
	userID1 := uuid.New()
	tetsChat := models2.Chat{
		ChatID: chatID,
	}
	testUser := models2.User{
		ID: userID1,
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareEchoCtx func(input models.AddUsersToChatRequest) (echo.Context, *httptest.ResponseRecorder)
		prepareRequest func() models.AddUsersToChatRequest
		prepare        func(f *fields, input models.AddUsersToChatRequest)
		checkResponse  func(recorder *httptest.ResponseRecorder) error
		wantErr        bool
	}{
		{
			name: "success",
			fields: fields{
				usecase: chat.NewMockChatUseCase(ctrl),
			},
			prepareRequest: func() models.AddUsersToChatRequest {
				input := models.AddUsersToChatRequest{
					ChatID:       &chatID,
					Participants: []uuid.UUID{userID1},
				}
				return input
			},
			prepareEchoCtx: func(input models.AddUsersToChatRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models.AddUsersToChatRequest) {
				f.usecase.EXPECT().AddUsersToChat(gomock.Any(), tetsChat, testUser).Return(models2.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatEchoHandler{
				usecase: tt.fields.usecase,
			}
			input := tt.prepareRequest()
			tt.prepare(&tt.fields, input)
			c, rec := tt.prepareEchoCtx(input)
			if err := ch.AddUsersToChat(c); (err != nil) != tt.wantErr {
				t.Errorf("AddUsersToChat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestChatEchoHandler_ChangeChatPhoto(t *testing.T) {
	type fields struct {
		usecase *chat.MockChatUseCase
	}
	type args struct {
		c echo.Context
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatID := uuid.New()
	testURL := "test"
	tetsChat := models2.Chat{
		ChatID: chatID,
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareEchoCtx func(input models.UpdateChatPhotoURLRequest) (echo.Context, *httptest.ResponseRecorder)
		prepareRequest func() models.UpdateChatPhotoURLRequest
		prepare        func(f *fields, input models.UpdateChatPhotoURLRequest)
		checkResponse  func(recorder *httptest.ResponseRecorder) error
		wantErr        bool
	}{
		{
			name: "success",
			fields: fields{
				usecase: chat.NewMockChatUseCase(ctrl),
			},
			prepareRequest: func() models.UpdateChatPhotoURLRequest {
				input := models.UpdateChatPhotoURLRequest{
					ChatID:   &chatID,
					PhotoURL: &testURL,
				}
				return input
			},
			prepareEchoCtx: func(input models.UpdateChatPhotoURLRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models.UpdateChatPhotoURLRequest) {
				f.usecase.EXPECT().UpdateChatPhotoURL(gomock.Any(), tetsChat, testURL).Return(models2.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatEchoHandler{
				usecase: tt.fields.usecase,
			}
			input := tt.prepareRequest()
			tt.prepare(&tt.fields, input)
			c, rec := tt.prepareEchoCtx(input)
			if err := ch.ChangeChatPhoto(c); (err != nil) != tt.wantErr {
				t.Errorf("ChangeChatPhoto() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestChatEchoHandler_PostNewChat(t *testing.T) {
	type fields struct {
		usecase *chat.MockChatUseCase
	}
	type args struct {
		c echo.Context
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()

	chatID := uuid.New()
	testURL := "test"
	testChat := models2.Chat{
		ChatID: chatID,
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareEchoCtx func(input models.CreateChatRequest) (echo.Context, *httptest.ResponseRecorder)
		prepareRequest func() models.CreateChatRequest
		prepare        func(f *fields, input models.CreateChatRequest)
		checkResponse  func(recorder *httptest.ResponseRecorder) error
		wantErr        bool
	}{
		{
			name: "success",
			fields: fields{
				usecase: chat.NewMockChatUseCase(ctrl),
			},
			prepareRequest: func() models.CreateChatRequest {
				input := models.CreateChatRequest{
					PhotoURL: &testURL,
					IssuerID: userID,
					Participants: []uuid.UUID{
						userID,
					},
				}
				return input
			},
			prepareEchoCtx: func(input models.CreateChatRequest) (echo.Context, *httptest.ResponseRecorder) {
				inputByte, _ := json.Marshal(&input)
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(inputByte)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", userID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields, input models.CreateChatRequest) {
				f.usecase.EXPECT().CreateChat(gomock.Any(), gomock.Cond(func(x any) bool {
					testChat := x.(models.CreateChatRequest)
					if len(testChat.Participants) != 1 {
						return false
					}
					if testChat.IssuerID == uuid.Nil {
						return false
					}
					if *testChat.PhotoURL == "" {
						return false
					}
					return true
				})).Return(testChat, models2.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusCreated {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatEchoHandler{
				usecase: tt.fields.usecase,
			}
			input := tt.prepareRequest()
			tt.prepare(&tt.fields, input)
			c, rec := tt.prepareEchoCtx(input)
			if err := ch.PostNewChat(c); (err != nil) != tt.wantErr {
				t.Errorf("PostNewChat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestChatEchoHandler_GetChat(t *testing.T) {
	type fields struct {
		usecase *chat.MockChatUseCase
	}
	type args struct {
		c echo.Context
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatID := uuid.New()
	testChat := models2.Chat{
		ChatID: chatID,
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareEchoCtx func() (echo.Context, *httptest.ResponseRecorder)
		prepare        func(f *fields)
		checkResponse  func(recorder *httptest.ResponseRecorder) error
		wantErr        bool
	}{
		{
			name: "success",
			fields: fields{
				usecase: chat.NewMockChatUseCase(ctrl),
			},
			args: args{},
			prepareEchoCtx: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.SetParamNames("id")
				testEchoCtx.SetParamValues(chatID.String())
				return testEchoCtx, rec
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					GetChat(gomock.Any(), models2.Chat{ChatID: chatID}).
					Return(testChat, models2.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatEchoHandler{
				usecase: tt.fields.usecase,
			}
			c, rec := tt.prepareEchoCtx()
			tt.prepare(&tt.fields)
			if err := ch.GetChat(c); (err != nil) != tt.wantErr {
				t.Errorf("GetChat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestChatEchoHandler_GetChatList(t *testing.T) {
	type fields struct {
		usecase *chat.MockChatUseCase
	}
	type args struct {
		c echo.Context
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()

	chatID := uuid.New()
	testChat := models2.ChatItem{
		ChatID: chatID,
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareEchoCtx func() (echo.Context, *httptest.ResponseRecorder)
		prepare        func(f *fields)
		checkResponse  func(recorder *httptest.ResponseRecorder) error
		wantErr        bool
	}{
		{
			name: "success",
			fields: fields{
				usecase: chat.NewMockChatUseCase(ctrl),
			},
			prepareEchoCtx: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.Set("user_id", userID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().
					GetChatList(gomock.Any(), models2.User{ID: userID}).
					Return([]models2.ChatItem{testChat}, models2.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatEchoHandler{
				usecase: tt.fields.usecase,
			}
			c, rec := tt.prepareEchoCtx()
			tt.prepare(&tt.fields)
			if err := ch.GetChatList(c); (err != nil) != tt.wantErr {
				t.Errorf("GetChatList() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestChatEchoHandler_GetChatMessages(t *testing.T) {
	type fields struct {
		usecase *chat.MockChatUseCase
	}
	type args struct {
		c echo.Context
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testOpts := models2.Opts{
		Page:  0,
		Limit: 1,
	}

	userID := uuid.New()

	chatID := uuid.New()
	testChat := models2.Chat{
		ChatID: chatID,
	}

	testMsg := models2.Message{
		MsgID:  uuid.New(),
		ChatID: chatID,
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareEchoCtx func() (echo.Context, *httptest.ResponseRecorder)
		prepare        func(f *fields)
		checkResponse  func(recorder *httptest.ResponseRecorder) error
		wantErr        bool
	}{
		{
			name: "success",
			fields: fields{
				usecase: chat.NewMockChatUseCase(ctrl),
			},
			prepareEchoCtx: func() (echo.Context, *httptest.ResponseRecorder) {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				testEchoCtx := e.NewContext(req, rec)
				testEchoCtx.SetParamNames("id")
				testEchoCtx.SetParamValues(chatID.String())
				testEchoCtx.QueryParams().Set("offset", "0")
				testEchoCtx.QueryParams().Set("limit", "1")
				testEchoCtx.Set("user_id", userID)
				return testEchoCtx, rec
			},
			prepare: func(f *fields) {
				f.usecase.EXPECT().GetChatMessages(gomock.Any(), testChat, testOpts).
					Return(models2.Messages{testMsg}, models2.OK)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) error {
				if recorder.Code != http.StatusOK {
					return fmt.Errorf("wrong status code")
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatEchoHandler{
				usecase: tt.fields.usecase,
			}
			c, rec := tt.prepareEchoCtx()
			tt.prepare(&tt.fields)
			if err := ch.GetChatMessages(c); (err != nil) != tt.wantErr {
				t.Errorf("GetChatMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkResponse(rec); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

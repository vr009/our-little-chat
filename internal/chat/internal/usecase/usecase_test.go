package usecase

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"our-little-chatik/internal/chat/internal/mocks/chat"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"reflect"
	"testing"
)

func TestChatUseCase_AddUsersToChat(t *testing.T) {
	type fields struct {
		repo  *chat.MockChatRepo
		queue *chat.MockQueueRepo
		users *chat.MockUserDataInteractor
	}
	type args struct {
		ctx   context.Context
		chat  models.Chat
		users []models.User
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testUserID1 := uuid.New()
	testUserID2 := uuid.New()
	testUserID3 := uuid.New()

	testUser3 := models.User{
		ID:       testUserID3,
		Name:     "test3",
		Nickname: "test3",
		Surname:  "test3",
	}
	testCtx := context.Background()

	testChat := models.Chat{
		ChatID: uuid.New(),
		Participants: []uuid.UUID{
			testUserID1,
			testUserID2,
		},
	}

	testNamedChat := models.Chat{
		ChatID: uuid.New(),
		Participants: []uuid.UUID{
			testUserID1,
			testUserID2,
		},
		Name: "chat",
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func(f *fields)
		status models.StatusCode
	}{
		{
			name: "success",
			fields: fields{
				repo:  chat.NewMockChatRepo(ctrl),
				queue: chat.NewMockQueueRepo(ctrl),
				users: chat.NewMockUserDataInteractor(ctrl),
			},
			args: args{
				ctx:   testCtx,
				chat:  testChat,
				users: []models.User{testUser3},
			},
			pre: func(f *fields) {
				f.repo.EXPECT().GetChat(testCtx, testChat).Return(testChat, models.OK)
				f.repo.EXPECT().AddUsersToChat(testCtx, testChat, gomock.Cond(func(x any) bool {
					chatNames := x.(map[string]string)
					if _, ok := chatNames[testUser3.ID.String()]; !ok {
						return false
					}
					return true
				}), testUser3)
			},
			status: models.OK,
		},
		{
			name: "success chat is named",
			fields: fields{
				repo:  chat.NewMockChatRepo(ctrl),
				queue: chat.NewMockQueueRepo(ctrl),
				users: chat.NewMockUserDataInteractor(ctrl),
			},
			args: args{
				ctx:   testCtx,
				chat:  testNamedChat,
				users: []models.User{testUser3},
			},
			pre: func(f *fields) {
				f.repo.EXPECT().GetChat(testCtx, testNamedChat).Return(testNamedChat, models.OK)
				f.repo.EXPECT().AddUsersToChat(testCtx, testNamedChat, gomock.Cond(func(x any) bool {
					chatNames := x.(map[string]string)
					if name, ok := chatNames[testUser3.ID.String()]; !ok && name != testNamedChat.Name {
						return false
					}
					return true
				}), testUser3)
			},
			status: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
				users: tt.fields.users,
			}
			tt.pre(&tt.fields)
			if status := ch.AddUsersToChat(tt.args.ctx, tt.args.chat, tt.args.users...); status != tt.status {
				t.Errorf("AddUsersToChat() error = %v, status %v", status, tt.status)
			}
		})
	}
}

func TestChatUseCase_CreateChat(t *testing.T) {
	type fields struct {
		repo  *chat.MockChatRepo
		queue *chat.MockQueueRepo
		users *chat.MockUserDataInteractor
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testName := "test1"
	testPhotoURL := "test"
	testUserID1 := uuid.New()
	testUserID2 := uuid.New()

	testUser1 := models.User{
		ID:       testUserID1,
		Name:     "test1",
		Nickname: "test1",
		Surname:  "test1",
	}
	testUser2 := models.User{
		ID:       testUserID2,
		Name:     "test2",
		Nickname: "test2",
		Surname:  "test2",
	}
	testCtx := context.Background()

	testChatRequest1 := models2.CreateChatRequest{
		Participants: []uuid.UUID{testUserID1, testUserID2},
		IssuerID:     testUserID1,
		Name:         &testName,
		PhotoURL:     &testPhotoURL,
	}

	testChatRequest2 := models2.CreateChatRequest{
		Participants: []uuid.UUID{testUserID2},
		IssuerID:     testUserID1,
		Name:         &testName,
		PhotoURL:     &testPhotoURL,
	}

	testChatRequest3 := models2.CreateChatRequest{
		Participants: []uuid.UUID{testUserID1},
		IssuerID:     testUserID1,
		Name:         &testName,
		PhotoURL:     &testPhotoURL,
	}

	type args struct {
		ctx     context.Context
		request models2.CreateChatRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func(f *fields)
		want   func(models.Chat) bool
		status models.StatusCode
	}{
		{
			name: "success",
			fields: fields{
				repo:  chat.NewMockChatRepo(ctrl),
				queue: chat.NewMockQueueRepo(ctrl),
				users: chat.NewMockUserDataInteractor(ctrl),
			},
			args: args{
				ctx:     testCtx,
				request: testChatRequest1,
			},
			pre: func(f *fields) {
				f.users.EXPECT().GetUser(models.User{ID: testUserID2}).Return(testUser2, models.OK)
				f.users.EXPECT().GetUser(models.User{ID: testUserID1}).Return(testUser1, models.OK)
				f.repo.EXPECT().CreateChat(testCtx, gomock.Cond(func(x any) bool {
					ch := x.(models.Chat)
					if ch.ChatID == uuid.Nil {
						return false
					}
					if ch.CreatedAt == 0 {
						return false
					}
					if len(ch.Participants) != 2 {
						return false
					}
					if ch.Name != testUser1.Name && ch.Name != testUser2.Name {
						return false
					}
					return true
				}), gomock.Cond(func(x any) bool {
					chatNames := x.(map[string]string)
					if chatNames[testUserID1.String()] != testUser2.Name {
						return false
					}
					if chatNames[testUserID2.String()] != testUser1.Name {
						return false
					}
					return true
				}))
			},
			want: func(m models.Chat) bool {
				if m.ChatID == uuid.Nil {
					return false
				}
				if len(m.Participants) != 2 {
					return false
				}
				if m.CreatedAt == 0 {
					return false
				}
				return true
			},
			status: models.OK,
		},
		{
			name: "success include self ",
			fields: fields{
				repo:  chat.NewMockChatRepo(ctrl),
				queue: chat.NewMockQueueRepo(ctrl),
				users: chat.NewMockUserDataInteractor(ctrl),
			},
			args: args{
				ctx:     testCtx,
				request: testChatRequest2,
			},
			pre: func(f *fields) {
				f.users.EXPECT().GetUser(models.User{ID: testUserID1}).Return(testUser1, models.OK)
				f.users.EXPECT().GetUser(models.User{ID: testUserID2}).Return(testUser2, models.OK)
				f.repo.EXPECT().CreateChat(testCtx, gomock.Cond(func(x any) bool {
					ch := x.(models.Chat)
					if ch.ChatID == uuid.Nil {
						return false
					}
					if ch.CreatedAt == 0 {
						return false
					}
					if len(ch.Participants) != 2 {
						return false
					}
					if ch.Name != testUser1.Name && ch.Name != testUser2.Name {
						return false
					}
					return true
				}), gomock.Cond(func(x any) bool {
					chatNames := x.(map[string]string)
					if chatNames[testUserID1.String()] != testUser2.Name {
						return false
					}
					if chatNames[testUserID2.String()] != testUser1.Name {
						return false
					}
					return true
				}))
			},
			status: models.OK,
			want: func(m models.Chat) bool {
				if m.ChatID == uuid.Nil {
					return false
				}
				if len(m.Participants) != 2 {
					return false
				}
				if m.CreatedAt == 0 {
					return false
				}
				return true
			},
		},
		{
			name: "success single user chat",
			fields: fields{
				repo:  chat.NewMockChatRepo(ctrl),
				queue: chat.NewMockQueueRepo(ctrl),
				users: chat.NewMockUserDataInteractor(ctrl),
			},
			args: args{
				ctx:     testCtx,
				request: testChatRequest3,
			},
			pre: func(f *fields) {
				f.users.EXPECT().GetUser(models.User{ID: testUserID1}).Return(testUser1, models.OK)
				f.repo.EXPECT().CreateChat(testCtx, gomock.Cond(func(x any) bool {
					ch := x.(models.Chat)
					if ch.ChatID == uuid.Nil {
						return false
					}
					if ch.CreatedAt == 0 {
						return false
					}
					if len(ch.Participants) != 1 {
						return false
					}
					if ch.Name != testUser1.Name && ch.Name != testUser2.Name {
						return false
					}
					return true
				}), gomock.Cond(func(x any) bool {
					chatNames := x.(map[string]string)
					if chatNames[testUserID1.String()] != testUser1.Name {
						return false
					}
					return true
				}))
			},
			status: models.OK,
			want: func(m models.Chat) bool {
				if m.ChatID == uuid.Nil {
					return false
				}
				if len(m.Participants) != 1 {
					return false
				}
				if m.CreatedAt == 0 {
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
				users: tt.fields.users,
			}
			tt.pre(&tt.fields)
			got, status := ch.CreateChat(tt.args.ctx, tt.args.request)
			if status != tt.status {
				t.Errorf("CreateChat() error = %v, status %v", status, tt.status)
				return
			}
			if !tt.want(got) {
				t.Errorf("failed check of the result")
			}
		})
	}
}

func TestChatUseCase_GetChatMessages(t *testing.T) {
	type fields struct {
		repo  *chat.MockChatRepo
		queue *chat.MockQueueRepo
		users *chat.MockUserDataInteractor
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testOpts1 := models.Opts{
		Limit: 3,
		Page:  0,
	}
	testOpts2 := models.Opts{
		Limit: 1,
		Page:  0,
	}

	testChat := models.Chat{
		ChatID: uuid.New(),
	}

	testMsg1 := models.Message{
		ChatID:    testChat.ChatID,
		MsgID:     uuid.New(),
		Payload:   "",
		CreatedAt: 10,
	}

	testMsg2 := models.Message{
		ChatID:    testChat.ChatID,
		MsgID:     uuid.New(),
		Payload:   "",
		CreatedAt: 9,
	}
	testMsg3 := models.Message{
		ChatID:    testChat.ChatID,
		MsgID:     uuid.New(),
		Payload:   "",
		CreatedAt: 8,
	}

	type args struct {
		ctx  context.Context
		chat models.Chat
		opts models.Opts
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func(f *fields)
		want   models.Messages
		status models.StatusCode
	}{
		{
			name: "fetch from both storages",
			fields: fields{
				repo:  chat.NewMockChatRepo(ctrl),
				queue: chat.NewMockQueueRepo(ctrl),
				users: chat.NewMockUserDataInteractor(ctrl),
			},
			args: args{
				ctx:  testCtx,
				chat: testChat,
				opts: testOpts1,
			},
			pre: func(f *fields) {
				f.queue.EXPECT().GetChatMessages(testChat, testOpts1).
					Return(models.Messages{testMsg1, testMsg2}, models.OK)
				f.repo.EXPECT().GetChatMessages(testCtx, testChat, testOpts2).
					Return(models.Messages{testMsg3}, models.OK)
			},
			want:   models.Messages{testMsg1, testMsg2, testMsg3},
			status: models.OK,
		},
		{
			name: "fetch only from queue",
			fields: fields{
				repo:  chat.NewMockChatRepo(ctrl),
				queue: chat.NewMockQueueRepo(ctrl),
				users: chat.NewMockUserDataInteractor(ctrl),
			},
			args: args{
				ctx:  testCtx,
				chat: testChat,
				opts: testOpts2,
			},
			pre: func(f *fields) {
				f.queue.EXPECT().GetChatMessages(testChat, testOpts2).
					Return(models.Messages{testMsg1}, models.OK)
			},
			want:   models.Messages{testMsg1},
			status: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
				users: tt.fields.users,
			}
			tt.pre(&tt.fields)
			got, status := ch.GetChatMessages(tt.args.ctx, tt.args.chat, tt.args.opts)
			if status != tt.status {
				t.Errorf("GetChatMessages() error = %v, status %v", status, tt.status)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChatMessages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

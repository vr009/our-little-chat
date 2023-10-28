package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"our-little-chatik/internal/models"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestPostgresRepo_FetchChatList(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		user models.User
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUserID := uuid.New()
	testChatID := uuid.New()

	testURL := "test_url"
	testName := "test"

	testMsg := models.Message{
		MsgID:     uuid.New(),
		SenderID:  uuid.New(),
		Payload:   "test",
		CreatedAt: int64(1),
	}

	testChat := models.Chat{
		ChatID:      testChatID,
		Name:        testName,
		PhotoURL:    testURL,
		LastMessage: testMsg,
	}

	columns := []string{
		"cp.chat_id",
		"cp.chat_name",
		"c.photo_url",
		"m.msg_id",
		"m.sender_id",
		"m.payload",
		"m.created_at",
	}

	tests := []struct {
		name   string
		pre    func()
		fields fields
		args   args
		want   []models.Chat
		status models.StatusCode
	}{
		{
			name: "successful fetch",
			fields: fields{
				pool: db,
			},
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(FetchChatListQuery)).
					WithArgs(testUserID).
					WillReturnRows(sqlmock.NewRows(columns).AddRow(testChatID,
						testName, testURL, testMsg.MsgID, testMsg.SenderID,
						testMsg.Payload, testMsg.CreatedAt))
			},
			args: args{
				user: models.User{
					ID: testUserID,
				},
			},
			want:   []models.Chat{testChat},
			status: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, status := pr.FetchChatList(context.Background(), tt.args.user)
			if status != tt.status {
				t.Errorf("FetchChatList() error = %v, wantErr %v", status, tt.status)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchChatList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_GetChat(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		chat models.Chat
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testChatID := uuid.New()
	testURL := "test_url"
	testName := "test"
	testTimestamp := time.Now().Unix()
	participant1 := uuid.New()
	participant2 := uuid.New()

	testMsg := models.Message{
		MsgID:     uuid.New(),
		SenderID:  uuid.New(),
		Payload:   "test",
		CreatedAt: int64(1),
	}

	passingTestChat := models.Chat{
		ChatID: testChatID,
	}
	expectedTestChat := models.Chat{
		ChatID:       testChatID,
		Name:         testName,
		PhotoURL:     testURL,
		CreatedAt:    testTimestamp,
		Participants: []uuid.UUID{participant1, participant2},
		LastMessage: models.Message{
			MsgID:     testMsg.MsgID,
			SenderID:  testMsg.SenderID,
			Payload:   testMsg.Payload,
			CreatedAt: testMsg.CreatedAt,
		},
	}

	columns := []string{
		"c.chat_id",
		"cp.chat_name",
		"c.photo_url",
		"c.created_at",
		"m.msg_id",
		"m.sender_id",
		"m.payload",
		"m.created_at",
	}

	pColumns := []string{
		"participant_id",
	}

	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		want   models.Chat
		status models.StatusCode
	}{
		{
			name: "Successful get",
			fields: fields{
				pool: db,
			},
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetChatInfoQuery)).
					WithArgs(expectedTestChat.ChatID).
					WillReturnRows(sqlmock.NewRows(columns).AddRow(testChatID,
						testName, testURL, testTimestamp, testMsg.MsgID, testMsg.SenderID,
						testMsg.Payload, testMsg.CreatedAt))
				mock.ExpectQuery(regexp.QuoteMeta(GetChatParticipantsQuery)).
					WithArgs(expectedTestChat.ChatID).
					WillReturnRows(sqlmock.NewRows(pColumns).AddRow(participant1).AddRow(participant2))
			},
			args: args{
				chat: passingTestChat,
			},
			want:   expectedTestChat,
			status: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, status := pr.GetChat(context.Background(), tt.args.chat)
			if status != tt.status {
				t.Errorf("GetChat() error = %v, wantErr %v", status, tt.status)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_GetChatMessages(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}
	type args struct {
		chat models.Chat
		opts models.Opts
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testMsgID := uuid.New()
	testUserID := uuid.New()
	testChatID := uuid.New()

	testPayload := "test_payload"
	testTimestamp := time.Now().Unix()

	testMsg := models.Message{
		ChatID:    testChatID,
		Payload:   testPayload,
		SenderID:  testUserID,
		MsgID:     testMsgID,
		CreatedAt: testTimestamp,
	}

	testChat := models.Chat{
		ChatID: testChatID,
	}

	columns := []string{
		"msg_id",
		"sender_id",
		"payload",
		"created_at",
	}

	tests := []struct {
		name   string
		pre    func()
		fields fields
		args   args
		want   models.Messages
		status models.StatusCode
	}{
		{
			name: "Successful",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetChatMessagesQuery)).
					WithArgs(testChatID, int64(0), int64(1)).
					WillReturnRows(sqlmock.NewRows(columns).AddRow(testMsgID,
						testUserID, testPayload, testTimestamp))
			},
			fields: fields{
				pool: db,
			},
			args: args{
				chat: testChat,
				opts: models.Opts{
					Page:  0,
					Limit: 1,
				},
			},
			want:   models.Messages{testMsg},
			status: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, status := pr.GetChatMessages(context.Background(), tt.args.chat, tt.args.opts)
			if status != tt.status {
				t.Errorf("GetChatMessages() error = %v, wantErr %v", status, tt.status)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChatMessages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_CreateChat(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUser1 := models.User{
		ID:   uuid.New(),
		Name: "test1",
	}
	testUser2 := models.User{
		ID:   uuid.New(),
		Name: "test2",
	}

	testChat := models.Chat{
		ChatID:   uuid.New(),
		PhotoURL: "test.png",
		Participants: []uuid.UUID{
			testUser1.ID,
			testUser2.ID,
		},
	}

	testCtx := context.Background()

	type args struct {
		ctx       context.Context
		chat      models.Chat
		chatNames map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		pre    func()
		args   args
		status models.StatusCode
	}{
		{
			name: "successful create",
			fields: fields{
				pool: db,
			},
			pre: func() {
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta(CreateChatParticipantsQuery)).
					WithArgs(testChat.ChatID, testUser1.ID, testUser2.Name).
					WillReturnError(nil).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(regexp.QuoteMeta(CreateChatParticipantsQuery)).
					WithArgs(testChat.ChatID, testUser2.ID, testUser1.Name).
					WillReturnError(nil).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(regexp.QuoteMeta(CreateChatQuery)).WithArgs(testChat.ChatID,
					testChat.PhotoURL, testChat.CreatedAt).
					WillReturnError(nil).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(nil)
			},
			args: args{
				ctx:  testCtx,
				chat: testChat,
				chatNames: map[string]string{
					testUser1.ID.String(): testUser2.Name,
					testUser2.ID.String(): testUser1.Name,
				},
			},
			status: models.OK,
		},
		{
			name: "impossible to create chat - already exists",
			fields: fields{
				pool: db,
			},
			pre: func() {
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta(CreateChatParticipantsQuery)).
					WithArgs(testChat.ChatID, testUser1.ID, testUser2.Name).
					WillReturnError(nil).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(regexp.QuoteMeta(CreateChatParticipantsQuery)).
					WithArgs(testChat.ChatID, testUser2.ID, testUser1.Name).
					WillReturnError(nil).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(regexp.QuoteMeta(CreateChatQuery)).WithArgs(testChat.ChatID,
					testChat.PhotoURL, testChat.CreatedAt).
					WillReturnError(nil).WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback().WillReturnError(nil)
			},
			args: args{
				ctx:  testCtx,
				chat: testChat,
				chatNames: map[string]string{
					testUser1.ID.String(): testUser2.Name,
					testUser2.ID.String(): testUser1.Name,
				},
			},
			status: models.InternalError,
		},
		{
			name: "fail to add a participant",
			fields: fields{
				pool: db,
			},
			pre: func() {
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta(CreateChatParticipantsQuery)).
					WithArgs(testChat.ChatID, testUser1.ID, testUser2.Name).
					WillReturnError(fmt.Errorf("")).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback().WillReturnError(nil)
			},
			args: args{
				ctx:  testCtx,
				chat: testChat,
				chatNames: map[string]string{
					testUser1.ID.String(): testUser2.Name,
					testUser2.ID.String(): testUser1.Name,
				},
			},
			status: models.InternalError,
		},
		{
			name: "fail to start transaction",
			fields: fields{
				pool: db,
			},
			pre: func() {
				mock.ExpectBegin().WillReturnError(fmt.Errorf(""))
			},
			args: args{
				ctx:  testCtx,
				chat: testChat,
				chatNames: map[string]string{
					testUser1.ID.String(): testUser2.Name,
					testUser2.ID.String(): testUser1.Name,
				},
			},
			status: models.InternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			if status := pr.CreateChat(tt.args.ctx, tt.args.chat,
				tt.args.chatNames); status != tt.status {
				t.Errorf("CreateChat() error = %v, wantErr %v", status, tt.status)
			}
		})
	}
}

func TestPostgresRepo_UpdateChatPhotoURL(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testChat := models.Chat{
		ChatID:   uuid.New(),
		PhotoURL: "test.png",
	}

	testCtx := context.Background()

	type args struct {
		ctx      context.Context
		chat     models.Chat
		photoURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func()
		status models.StatusCode
	}{
		{
			name: "success",
			fields: fields{
				pool: db,
			},
			args: args{
				ctx:      testCtx,
				chat:     testChat,
				photoURL: testChat.PhotoURL,
			},
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(UpdatePhotoURLQuery)).
					WithArgs(testChat.PhotoURL, testChat.ChatID).
					WillReturnError(nil).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			status: models.OK,
		},
		{
			name: "fail",
			fields: fields{
				pool: db,
			},
			args: args{
				ctx:      testCtx,
				chat:     testChat,
				photoURL: testChat.PhotoURL,
			},
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(UpdatePhotoURLQuery)).
					WithArgs(testChat.PhotoURL, testChat.ChatID).
					WillReturnError(nil).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			status: models.InternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			if status := pr.UpdateChatPhotoURL(tt.args.ctx, tt.args.chat, tt.args.photoURL); status != tt.status {
				t.Errorf("UpdateChatPhotoURL() error = %v, wantErr %v", status, tt.status)
			}
		})
	}
}

func TestPostgresRepo_AddUsersToChat(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUser1 := models.User{
		ID:   uuid.New(),
		Name: "test1",
	}
	testUser2 := models.User{
		ID:   uuid.New(),
		Name: "test2",
	}

	testChat := models.Chat{
		ChatID:   uuid.New(),
		PhotoURL: "test.png",
		Participants: []uuid.UUID{
			testUser1.ID,
			testUser2.ID,
		},
	}

	testCtx := context.Background()

	type args struct {
		ctx   context.Context
		chat  models.Chat
		users []models.User
	}
	testChatNames := map[string]string{
		testUser1.ID.String(): "group",
		testUser2.ID.String(): "group2",
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func()
		status models.StatusCode
	}{
		{
			name: "success",
			fields: fields{
				pool: db,
			},
			args: args{
				ctx:   testCtx,
				chat:  testChat,
				users: []models.User{testUser1, testUser2},
			},
			pre: func() {
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta(CreateChatParticipantsQuery)).
					WithArgs(testChat.ChatID, testUser1.ID, testChatNames[testUser1.ID.String()]).
					WillReturnError(nil).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(regexp.QuoteMeta(CreateChatParticipantsQuery)).
					WithArgs(testChat.ChatID, testUser2.ID, testChatNames[testUser2.ID.String()]).
					WillReturnError(nil).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(nil)
			},
			status: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			if status := pr.AddUsersToChat(tt.args.ctx, tt.args.chat, testChatNames,
				tt.args.users...); status != tt.status {
				t.Errorf("AddUsersToChat() error = %v, wantErr %v", status, tt.status)
			}
		})
	}
}

func TestPostgresRepo_RemoveUserFromChat(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testUser1 := models.User{
		ID:   uuid.New(),
		Name: "test1",
	}
	testUser2 := models.User{
		ID:   uuid.New(),
		Name: "test2",
	}

	testChat := models.Chat{
		ChatID:   uuid.New(),
		PhotoURL: "test.png",
		Participants: []uuid.UUID{
			testUser1.ID,
			testUser2.ID,
		},
	}

	testCtx := context.Background()

	type args struct {
		ctx   context.Context
		chat  models.Chat
		users []models.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func()
		status models.StatusCode
	}{
		{
			name: "success",
			fields: fields{
				pool: db,
			},
			args: args{
				ctx:   testCtx,
				chat:  testChat,
				users: []models.User{testUser1},
			},
			pre: func() {
				mock.ExpectBegin().WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta(RemoveUserFromChatQuery)).
					WithArgs(testUser1.ID, testChat.ChatID).
					WillReturnError(nil).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(nil)
			},
			status: models.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			if status := pr.RemoveUserFromChat(tt.args.ctx, tt.args.chat, tt.args.users...); status != tt.status {
				t.Errorf("RemoveUserFromChat() error = %v, wantErr %v", status, tt.status)
			}
		})
	}
}

func TestPostgresRepo_DeleteChat(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testChat := models.Chat{
		ChatID:   uuid.New(),
		PhotoURL: "test.png",
	}

	testCtx := context.Background()

	type args struct {
		ctx  context.Context
		chat models.Chat
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func()
		status models.StatusCode
	}{
		{
			name: "success",
			fields: fields{
				pool: db,
			},
			args: args{
				ctx:  testCtx,
				chat: testChat,
			},
			status: models.Deleted,
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(DeleteChatQuery)).
					WithArgs(testChat.ChatID).
					WillReturnError(nil).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			if status := pr.DeleteChat(tt.args.ctx, tt.args.chat); status != tt.status {
				t.Errorf("DeleteChat() error = %v, wantErr %v", status, tt.status)
			}
		})
	}
}

func TestPostgresRepo_DeleteMessage(t *testing.T) {
	type fields struct {
		pool *sql.DB
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testMsg := models.Message{
		MsgID: uuid.New(),
	}
	testCtx := context.Background()

	type args struct {
		ctx     context.Context
		message models.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		pre    func()
		status models.StatusCode
	}{
		{
			name: "success",
			fields: fields{
				pool: db,
			},
			args: args{
				ctx:     testCtx,
				message: testMsg,
			},
			pre: func() {
				mock.ExpectExec(regexp.QuoteMeta(DeleteMessageQuery)).
					WithArgs(testMsg.MsgID).
					WillReturnError(nil).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			status: models.Deleted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			if status := pr.DeleteMessage(tt.args.ctx, tt.args.message); status != tt.status {
				t.Errorf("DeleteMessage() error = %v, wantErr %v", err, tt.status)
			}
		})
	}
}

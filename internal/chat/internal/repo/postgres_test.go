package repo

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"our-little-chatik/internal/models"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestPostgresRepo_FetchChatList(t *testing.T) {
	type fields struct {
		pool pgx.Tx
	}
	type args struct {
		user models.User
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testUserID := uuid.New()
	testChatID := uuid.New()

	testURL := "test_url"
	testName := "test"

	testChatItem := models.ChatItem{
		ChatID:   testChatID,
		Name:     testName,
		PhotoURL: testURL,
	}

	columns := []string{
		"cp.chat_id",
		"cp.chat_name",
		"c.photo_url",
	}

	tests := []struct {
		name    string
		pre     func()
		fields  fields
		args    args
		want    []models.ChatItem
		wantErr bool
	}{
		{
			name: "successful fetch",
			fields: fields{
				pool: mock,
			},
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(FetchChatListQuery)).
					WithArgs(testUserID).
					WillReturnRows(pgxmock.NewRows(columns).AddRow(testChatID,
						testName, testURL))
			},
			args: args{
				user: models.User{
					UserID: testUserID,
				},
			},
			want:    []models.ChatItem{testChatItem},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, err := pr.FetchChatList(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchChatList() error = %v, wantErr %v", err, tt.wantErr)
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
		pool pgx.Tx
	}
	type args struct {
		chat models.Chat
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	testChatID := uuid.New()
	testURL := "test_url"
	testName := "test"
	testTimestamp := time.Now().Unix()
	participant1 := uuid.New()
	participant2 := uuid.New()

	passingTestChat := models.Chat{
		ChatID: testChatID,
	}
	expectedTestChat := models.Chat{
		ChatID:       testChatID,
		Name:         testName,
		PhotoURL:     testURL,
		CreatedAt:    testTimestamp,
		Participants: []uuid.UUID{participant1, participant2},
	}

	columns := []string{
		"c.chat_id",
		"cp.chat_name",
		"c.photo_url",
		"c.created_at",
	}

	pColumns := []string{
		"participant_id",
	}

	tests := []struct {
		name    string
		fields  fields
		pre     func()
		args    args
		want    models.Chat
		wantErr bool
	}{
		{
			name: "Successful get",
			fields: fields{
				pool: mock,
			},
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetChatInfoQuery)).
					WithArgs(expectedTestChat.ChatID).
					WillReturnRows(pgxmock.NewRows(columns).AddRow(testChatID,
						testName, testURL, testTimestamp))
				mock.ExpectQuery(regexp.QuoteMeta(GetChatParticipantsQuery)).
					WithArgs(expectedTestChat.ChatID).
					WillReturnRows(pgxmock.NewRows(pColumns).AddRow(participant1).AddRow(participant2))
			},
			args: args{
				chat: passingTestChat,
			},
			want:    expectedTestChat,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, err := pr.GetChat(tt.args.chat)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChat() error = %v, wantErr %v", err, tt.wantErr)
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
		pool pgx.Tx
	}
	type args struct {
		chat models.Chat
		opts models.Opts
	}

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

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
		name    string
		pre     func()
		fields  fields
		args    args
		want    models.Messages
		wantErr bool
	}{
		{
			name: "Successful",
			pre: func() {
				mock.ExpectQuery(regexp.QuoteMeta(GetChatMessagesQuery)).
					WithArgs(testChatID, int64(0), int64(1)).
					WillReturnRows(pgxmock.NewRows(columns).AddRow(testMsgID,
						testUserID, testPayload, testTimestamp))
			},
			fields: fields{
				pool: mock,
			},
			args: args{
				chat: testChat,
				opts: models.Opts{
					Page:  0,
					Limit: 1,
				},
			},
			want:    models.Messages{testMsg},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := PostgresRepo{
				pool: tt.fields.pool,
			}
			tt.pre()
			got, err := pr.GetChatMessages(tt.args.chat, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChatMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChatMessages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

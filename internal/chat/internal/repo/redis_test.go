package repo

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"our-little-chatik/internal/models"
	"reflect"
	"testing"
	"time"
)

func TestRedisRepo_GetChatMessages(t *testing.T) {
	type fields struct {
		cl *redis.Client
	}
	type args struct {
		chat models.Chat
	}

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

	db, mock := redismock.NewClientMock()
	key := fmt.Sprintf("%s_%s", testChat.ChatID.String(), testMsgID.String())

	bData, _ := json.Marshal(testMsg)
	mock.ExpectKeys(testChat.ChatID.String() + "*").SetVal([]string{key})
	mock.ExpectMGet(key).SetVal([]interface{}{string(bData)})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Messages
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				cl: db,
			},
			args: args{
				chat: testChat,
			},
			want:    models.Messages{testMsg},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RedisRepo{
				cl: tt.fields.cl,
			}
			got, err := r.GetChatMessages(tt.args.chat, models.Opts{Limit: 10, Page: 0})
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

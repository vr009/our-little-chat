package repo

import (
	"encoding/json"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"our-little-chatik/internal/models"
	"reflect"
	"testing"
)

func TestRedisRepo_FetchAllMessages(t *testing.T) {
	type fields struct {
		cl *redis.Client
	}

	db, mock := redismock.NewClientMock()

	testMsg := models.Message{
		MsgID: uuid.New(),
	}
	keys := []string{
		"test",
	}

	testMsgByte, _ := json.Marshal(testMsg)

	tests := []struct {
		name    string
		fields  fields
		want    []models.Message
		pre     func()
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				cl: db,
			},
			want: []models.Message{testMsg},
			pre: func() {
				mock.ExpectKeys("*").SetVal(keys)
				mock.ExpectMGet(keys...).SetVal([]interface{}{string(testMsgByte)})
				mock.ExpectDel(keys...).SetVal(1)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RedisRepo{
				cl: tt.fields.cl,
			}
			tt.pre()
			got, err := r.FetchAllMessages()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchAllMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchAllMessages() got = %v, want %v", got, tt.want)
			}
		})
	}
}

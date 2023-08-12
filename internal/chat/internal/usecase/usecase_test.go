package usecase

import (
	"github.com/google/uuid"
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/chat/internal/repo"
	"our-little-chatik/internal/models"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestChatUseCase_CreateChat(t *testing.T) {
	type fields struct {
		repo  internal.ChatRepo
		queue internal.QueueRepo
	}
	type args struct {
		chat models.Chat
	}

	rMock := repo.RedisMock{}

	pMock := repo.PostgresMock{}

	testChat := models.Chat{
		Name: "test",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Chat
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{repo: pMock, queue: rMock},
			args: args{
				chat: testChat,
			},
			want:    testChat,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
			}
			got, err := ch.CreateChat(tt.args.chat)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("CreateChat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatUseCase_GetChatMessages(t *testing.T) {
	type fields struct {
		repo  internal.ChatRepo
		queue internal.QueueRepo
	}
	type args struct {
		chat models.Chat
		opts models.Opts
	}

	chatID := uuid.New()
	senderID := uuid.New()
	timestamp := time.Now().Unix()

	msg1 := models.Message{
		MsgID:     uuid.New(),
		Payload:   "test",
		SenderID:  senderID,
		ChatID:    chatID,
		CreatedAt: timestamp - 500,
	}
	msg2 := models.Message{
		MsgID:     uuid.New(),
		Payload:   "test1",
		SenderID:  senderID,
		ChatID:    chatID,
		CreatedAt: timestamp - 100,
	}

	msg3 := models.Message{
		MsgID:     uuid.New(),
		Payload:   "test",
		SenderID:  senderID,
		ChatID:    chatID,
		CreatedAt: timestamp - 1500,
	}
	msg4 := models.Message{
		MsgID:     uuid.New(),
		Payload:   "test1",
		SenderID:  senderID,
		ChatID:    chatID,
		CreatedAt: timestamp - 1200,
	}

	rMock := repo.RedisMock{
		ID: chatID,
		Msgs: models.Messages{
			msg1, msg2,
		},
	}

	pMock := repo.PostgresMock{
		Chat: models.Chat{
			ChatID: chatID,
		},
		Msgs: models.Messages{
			msg3, msg4,
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Messages
		wantErr bool
	}{
		{
			name: "Successful get",
			fields: fields{
				queue: rMock,
				repo:  pMock,
			},
			args:    args{chat: models.Chat{ChatID: chatID}, opts: models.Opts{Limit: 10, Page: 0}},
			wantErr: false,
			want: models.Messages{
				msg2,
				msg1,
				msg4,
				msg3,
			},
		},
		{
			name: "Successful get no data from redis",
			fields: fields{
				queue: repo.RedisMock{
					ID:   chatID,
					Msgs: models.Messages{},
				},
				repo: pMock,
			},
			args:    args{chat: models.Chat{ChatID: chatID}, opts: models.Opts{Limit: 10, Page: 0}},
			wantErr: false,
			want: models.Messages{
				msg4,
				msg3,
			},
		},
		{
			name: "Successful get no data from postgres",
			fields: fields{
				queue: rMock,
				repo: repo.PostgresMock{
					Chat: models.Chat{
						ChatID: chatID,
					},
					Msgs: models.Messages{},
				},
			},
			args:    args{chat: models.Chat{ChatID: chatID}, opts: models.Opts{Limit: 10, Page: 0}},
			wantErr: false,
			want: models.Messages{
				msg2,
				msg1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
			}
			got, err := ch.GetChatMessages(tt.args.chat, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChatMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChatMessages() got = %v, want %v", got, tt.want)
			}
			if !sort.IsSorted(got) {
				t.Errorf("result slice of messages is not sorted")
			}
		})
	}
}

func TestChatUseCase_GetChat(t *testing.T) {
	type fields struct {
		repo  internal.ChatRepo
		queue internal.QueueRepo
	}
	type args struct {
		chat models.Chat
	}
	rMock := repo.RedisMock{}

	pMock := repo.PostgresMock{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Chat
		wantErr bool
	}{
		{
			fields: fields{
				repo:  pMock,
				queue: rMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
			}
			got, err := ch.GetChat(tt.args.chat)
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

func TestChatUseCase_GetChatList(t *testing.T) {
	type fields struct {
		repo  internal.ChatRepo
		queue internal.QueueRepo
	}
	type args struct {
		user models.User
	}
	rMock := repo.RedisMock{}

	pMock := repo.PostgresMock{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.ChatItem
		wantErr bool
	}{
		{
			fields: fields{
				repo:  pMock,
				queue: rMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
			}
			got, err := ch.GetChatList(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChatList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChatList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatUseCase_UpdateChat(t *testing.T) {
	type fields struct {
		repo  internal.ChatRepo
		queue internal.QueueRepo
	}
	type args struct {
		chat       models.Chat
		updateOpts models2.UpdateOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &ChatUseCase{
				repo:  tt.fields.repo,
				queue: tt.fields.queue,
			}
			if err := ch.UpdateChat(tt.args.chat, tt.args.updateOpts); (err != nil) != tt.wantErr {
				t.Errorf("UpdateChat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

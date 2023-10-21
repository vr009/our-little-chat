package delivery

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"our-little-chatik/internal/flusher/internal/mocks/flusher"
	"our-little-chatik/internal/models"
	"testing"
	"time"
)

func TestFlusherD_Work(t *testing.T) {
	type fields struct {
		queueRepo      *flusher.MockQueueRepo
		persistantRepo *flusher.MockPersistantRepo
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx    context.Context
		period time.Duration
	}
	testMsg := models.Message{
		MsgID: uuid.New(),
	}

	testCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name    string
		prepare func(f *fields)
		fields  fields
		args    args
	}{
		{
			name: "success",
			prepare: func(f *fields) {
				f.queueRepo.EXPECT().FetchAllMessages().
					Return([]models.Message{testMsg}, nil)
				f.persistantRepo.EXPECT().PersistAllMessages([]models.Message{testMsg}).
					Return(nil)
			},
			args: args{
				ctx:    testCtx,
				period: time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &FlusherD{
				queueRepo:      tt.fields.queueRepo,
				persistantRepo: tt.fields.persistantRepo,
			}
			go func() {
				time.Sleep(time.Millisecond * 10)
				cancel()
			}()
			d.Work(tt.args.ctx, tt.args.period)
		})
	}
}

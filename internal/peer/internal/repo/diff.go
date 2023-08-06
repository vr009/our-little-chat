package repo

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"log"
	"our-little-chatik/internal/models"
)

type DiffRepository struct {
	cl *redis.Client
}

func NewDiffRepository(cl *redis.Client) *DiffRepository {
	return &DiffRepository{
		cl: cl,
	}
}

func (r *DiffRepository) StartSubscriber(ctx context.Context,
	messageChan chan models.Message, chatChannels []string) {
	/*
		this goroutine exits when the application shuts down. When the pusub connection is closed,
		the channel range loop terminates, hence terminating the goroutine
	*/
	go func() {
		log.Println("starting subscriber...")
		sub := r.cl.Subscribe(chatChannels...)
		messages := sub.Channel()
		for message := range messages {
			msg, err := parseMessage(message.Payload)
			if err != nil {
				glog.Error(err)
				continue
			}
			messageChan <- *msg
		}
		select {
		case <-ctx.Done():
			err := sub.Unsubscribe(chatChannels...)
			if err != nil {
				glog.Error(err)
			}
			return
		default:
		}
	}()
}

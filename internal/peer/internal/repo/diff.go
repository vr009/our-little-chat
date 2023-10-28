package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/models"
	models2 "our-little-chatik/internal/peer/internal/models"
)

type DiffRepository struct {
	cl *redis.Client
}

func NewDiffRepository(cl *redis.Client) *DiffRepository {
	return &DiffRepository{
		cl: cl,
	}
}

func (r *DiffRepository) SubscribeToChats(ctx context.Context,
	chats []models.Chat) (chan models.Message, error) {
	chatChannels := make([]string, 0)
	for _, chat := range chats {
		chatChannels = append(chatChannels, fmt.Sprintf(models2.CommonFormat, "chat", chat.ChatID.String()))
	}
	sub := r.cl.PSubscribe(chatChannels...)
	userMsgChan := make(chan models.Message)
	msgChan := sub.Channel()
	go func() {
		for {
			select {
			case redisMsg := <-msgChan:
				msg := models.Message{}
				bMsg := redisMsg.Payload
				err := json.Unmarshal([]byte(bMsg), &msg)
				if err != nil {
					slog.Error(err.Error())
					continue
				}
				userMsgChan <- msg
			case <-ctx.Done():
				slog.Warn("finish by context")
			}
		}
	}()
	return userMsgChan, nil
}

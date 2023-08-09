package repo

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/redis/go-redis/v9"
	"our-little-chatik/internal/models"
)

type RedisRepo struct {
	cl *redis.Client
}

func NewRedisRepo(cl *redis.Client) *RedisRepo {
	return &RedisRepo{
		cl: cl,
	}
}

func (r RedisRepo) GetFreshMessagesFromChat(chat models.Chat) (models.Messages, error) {
	keys, err := r.cl.Keys(context.Background(), chat.ChatID.String()+"*").Result()
	if err != nil {
		return nil, err
	}
	msgList := make(models.Messages, 0)
	values, err := r.cl.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}

	for _, val := range values {
		msg := models.Message{}
		err := json.Unmarshal([]byte(val.(string)), &msg)
		if err != nil {
			glog.Error(err)
			continue
		}
		msgList = append(msgList, msg)
	}
	return msgList, nil
}

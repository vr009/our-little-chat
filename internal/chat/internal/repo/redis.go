package repo

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	models2 "our-little-chatik/internal/chat/models"
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

func (r RedisRepo) GetFreshMessagesFromChat(chat models2.Chat) ([]models.Message, error) {
	keys, err := r.cl.Keys(chat.ChatID.String() + "*").Result()
	if err != nil {
		return nil, err
	}
	msgList := make([]models.Message, 0)
	values, err := r.cl.MGet(keys...).Result()
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

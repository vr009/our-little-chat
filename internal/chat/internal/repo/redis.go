package repo

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/models"
	"sort"
)

type RedisRepo struct {
	cl *redis.Client
}

func NewRedisRepo(cl *redis.Client) *RedisRepo {
	return &RedisRepo{
		cl: cl,
	}
}

func (r RedisRepo) GetChatMessages(chat models.Chat,
	opts models.Opts) (models.Messages, models.StatusCode) {
	keys, err := r.cl.Keys(context.Background(), chat.ChatID.String()+"*").Result()
	if err != nil {
		return nil, models.NotFound
	}
	msgList := make(models.Messages, 0)
	values, err := r.cl.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, models.NotFound
	}

	for _, val := range values {
		msg := models.Message{}
		err := json.Unmarshal([]byte(val.(string)), &msg)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		msgList = append(msgList, msg)
	}
	sort.Sort(msgList)

	limit := int(opts.Limit)
	page := int(opts.Page)
	firstElemIdx := limit * page
	lastElemIdx := limit*page + limit
	if lastElemIdx > len(msgList) {
		lastElemIdx = len(msgList)
	}
	if len(msgList) <= firstElemIdx {
		return models.Messages{}, models.NotFound
	} else {
		return msgList[firstElemIdx:lastElemIdx], models.OK
	}
}

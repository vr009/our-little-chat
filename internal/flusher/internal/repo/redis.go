package repo

import (
	"context"
	"encoding/json"
	"github.com/prometheus/common/log"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/models"
	"strings"
)

type RedisRepo struct {
	cl *redis.Client
}

func NewRedisRepo(cl *redis.Client) *RedisRepo {
	return &RedisRepo{
		cl: cl,
	}
}

func (r RedisRepo) FetchAllMessages() ([]models.Message, error) {
	keys, err := r.cl.Keys(context.Background(), "*").Result()
	if err != nil {
		return nil, err
	}
	filteredKeys := make([]string, 0)
	for _, key := range keys {
		if !strings.Contains(key, "last-inserted") {
			filteredKeys = append(filteredKeys, key)
		}
	}

	values, err := r.cl.MGet(context.Background(), filteredKeys...).Result()
	if err != nil {
		return nil, err
	}

	num, err := r.cl.Del(context.Background(), filteredKeys...).Result()
	if err != nil {
		log.Errorf("failed to delete keys in db: %v", err)
	}

	log.Infof("removed %d records", num)

	// We merge result of both requests in one slice. The same order of keys and values
	// in both slices is ensured.
	messages := make([]models.Message, 0)
	for i := range values {
		var msg models.Message
		valStr := values[i].(string)
		err := json.Unmarshal([]byte(valStr), &msg)
		if err != nil {
			slog.Warn("failed to cast a value from redis to models.Message")
			continue
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r RedisRepo) FetchAllLastMessagesOfChats() ([]models.Message, error) {
	keys, err := r.cl.Keys(context.Background(), "*:last-inserted").Result()
	if err != nil {
		return nil, err
	}

	values, err := r.cl.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}
	num, err := r.cl.Del(context.Background(), keys...).Result()
	if err != nil {
		log.Errorf("failed to delete keys in db: %v", err)
	}

	log.Infof("removed %d records", num)

	// We merge result of both requests in one slice. The same order of keys and values
	// in both slices is ensured.
	messages := make([]models.Message, 0)
	for i := range values {
		var msg models.Message
		valStr := values[i].(string)
		err := json.Unmarshal([]byte(valStr), &msg)
		if err != nil {
			slog.Warn("failed to cast a value from redis to models.Message")
			continue
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

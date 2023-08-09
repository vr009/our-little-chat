package repo

import (
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"github.com/prometheus/common/log"
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

func (r RedisRepo) FetchAllMessages() ([]models.Message, error) {
	keys, err := r.cl.Keys("*").Result()
	if err != nil {
		return nil, err
	}

	values, err := r.cl.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	num, err := r.cl.Del(keys...).Result()
	if err != nil {
		log.Errorf("failed to delete keys in db: %v", err)
	}

	log.Infof("removed %d records", num)

	// We merge result of both requests in one slice. The same order of keys and values
	// in both slices is ensured.
	messages := make([]models.Message, 0)
	for i := range values {
		if msg, ok := values[i].(models.Message); ok {
			messages = append(messages, msg)
		} else {
			glog.Warning("failed to cast a value from redis to models.Message")
		}
	}
	return messages, nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"our-little-chatik/internal/flusher/internal/delivery"
	"our-little-chatik/internal/flusher/internal/repo"
	"strconv"

	"github.com/golang/glog"
	"github.com/jackc/pgx/v5"
)

type PostgresConfig struct {
	URI      string
	Username string
	Password string
}

type PeerDBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type AppConfig struct {
	Port   string
	DB     PostgresConfig
	Redis  PeerDBConfig
	Period int
}

func GetConnectionString() (string, error) {
	key, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		panic("empty redis host")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		panic("empty redis port")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		panic("empty redis password")
	}
	flusherPort := os.Getenv("FLUSHER_PORT")
	if flusherPort == "" {
		panic("empty flusher port")
	}
	flusherPeriod := os.Getenv("FLUSHER_PERIOD")
	if flusherPeriod == "" {
		panic("empty flusher period")
	}

	period, err := strconv.Atoi(flusherPeriod)
	if err != nil {
		panic(err.Error())
	}

	appConfig := AppConfig{}
	appConfig.Port = flusherPort
	appConfig.Redis.Host = redisHost
	appConfig.Redis.Port = redisPort
	appConfig.Redis.Password = redisPassword
	appConfig.Period = period

	glog.V(2)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     appConfig.Redis.Host + ":" + appConfig.Redis.Port,
		Password: appConfig.Redis.Password,
	})

	ctx := context.Background()
	connStr, err := GetConnectionString()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		panic(err)
	}
	peristRepo := repo.NewPostgresRepo(conn)
	queueRepo := repo.NewRedisRepo(redisClient)

	daemon := delivery.NewFlusherD(queueRepo, peristRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println(appConfig.Period)
	daemon.Work(ctx, appConfig.Period)
}

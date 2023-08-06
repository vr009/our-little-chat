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

	"github.com/golang/glog"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
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
	Port   int
	DB     PostgresConfig
	PeerDB PeerDBConfig
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
	configPath := os.Getenv("FLUSHER_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("flusher-config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file")
	}

	appConfig := AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Fatal(err)
	}
	glog.V(2)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     appConfig.PeerDB.Host + ":" + appConfig.PeerDB.Port,
		Password: appConfig.PeerDB.Password,
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

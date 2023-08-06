package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"strconv"

	"our-little-chatik/internal/chat/internal/delivery"
	"our-little-chatik/internal/chat/internal/repo"
	"our-little-chatik/internal/chat/internal/usecase"
	"our-little-chatik/internal/chat/middleware"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type AppConfig struct {
	Port  int
	Redis RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func GetConnectionString() (string, error) {
	key, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func main() {
	configPath := os.Getenv("CHAT_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("chat-config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file ", err)
	}

	appConfig := AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Fatal(err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     appConfig.Redis.Host + ":" + appConfig.Redis.Port,
		Password: appConfig.Redis.Password,
	})

	ctx := context.Background()
	connStr, err := GetConnectionString()
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	repop := repo.NewPostgresRepo(pool)
	repoTT := repo.NewRedisRepo(redisClient)
	uc := usecase.NewChatUseCase(repop, repoTT)

	handler := delivery.NewChatHandler(uc)

	r := mux.NewRouter()

	// Getting chat messages
	r.HandleFunc("/api/v1/conv", handler.GetChatMessages).Methods("GET")
	// Getting the list of users chats
	r.HandleFunc("/api/v1/list", handler.GetChatList).Methods("GET")
	// Creating a new chat
	r.HandleFunc("/api/v1/new", handler.PostNewChat).Methods("POST")

	srv := &http.Server{Handler: middleware.Logger(r),
		Addr: ":" + strconv.Itoa(appConfig.Port)}

	log.Printf("Listening port: %d", appConfig.Port)
	log.Printf("addres to query: %s", "http://localhost:"+strconv.Itoa(appConfig.Port)+"/api/v1/")
	log.Fatal(srv.ListenAndServe())
}

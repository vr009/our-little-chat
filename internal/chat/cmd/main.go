package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strconv"

	"our-little-chatik/internal/chat/internal/delivery"
	"our-little-chatik/internal/chat/internal/repo"
	"our-little-chatik/internal/chat/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

type AppConfig struct {
	Port  int
	Redis RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     string
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
	var err error
	appConfig := AppConfig{}
	appConfig.Port, err = strconv.Atoi(os.Getenv("CHAT_PORT"))
	if err != nil {
		panic(err.Error())
	}
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

	appConfig.Redis.Port = redisPort
	appConfig.Redis.Host = redisHost
	appConfig.Redis.Password = redisPassword

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
	repoRed := repo.NewRedisRepo(redisClient)
	uc := usecase.NewChatUseCase(repop, repoRed)

	handler := delivery.NewChatHandler(uc)

	r := mux.NewRouter()

	// Getting chat info
	r.HandleFunc("/api/v1/chat", handler.GetChat).Methods("GET")
	// Getting chat messages
	r.HandleFunc("/api/v1/conv", handler.GetChatMessages).Methods("GET")
	// Getting the list of users chats
	r.HandleFunc("/api/v1/list", handler.GetChatList).Methods("GET")
	// Creating a new chat
	r.HandleFunc("/api/v1/new", handler.PostNewChat).Methods("POST")
	// Update photo url of the chat
	r.HandleFunc("/api/v1/chat/photo", handler.ChangeChatPhoto).Methods("POST")
	// Add users to chat
	r.HandleFunc("/api/v1/chat/users", handler.AddUsersToChat).Methods("POST")

	srv := &http.Server{Handler: r,
		Addr: ":" + strconv.Itoa(appConfig.Port)}

	log.Printf("Listening port: %d", appConfig.Port)
	log.Printf("addres to query: %s", "http://localhost:"+strconv.Itoa(appConfig.Port)+"/api/v1/")
	log.Fatal(srv.ListenAndServe())
}

package main

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
	"our-little-chatik/internal/peer/internal/delivery"
	"our-little-chatik/internal/peer/internal/repo"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type AppConfig struct {
	Port  string
	Redis DBConfig
}

func main() {
	appConfig := AppConfig{}
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
	peerPort := os.Getenv("PEER_PORT")
	if peerPort == "" {
		panic("empty peer port")
	}

	appConfig.Port = peerPort
	appConfig.Redis.Port = redisPort
	appConfig.Redis.Host = redisHost
	appConfig.Redis.Password = redisPassword

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     appConfig.Redis.Host + ":" + appConfig.Redis.Port,
		Password: appConfig.Redis.Password,
	})
	err := redisClient.Ping().Err()
	if err != nil {
		panic(err)
	}
	peerRepo := repo.NewPeerRepository(redisClient)
	peerHandler := delivery.NewPeerHandler(peerRepo, peerRepo)

	diffRepo := repo.NewDiffRepository(redisClient)

	diffHandler := delivery.NewDiffHandler(peerRepo, diffRepo)

	r := mux.NewRouter()

	r.HandleFunc("/ws/chat", peerHandler.ConnectToChat)
	r.HandleFunc("/ws/diff", diffHandler.ConnectToDiff)

	slog.Info("service started", "port", appConfig.Port)

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + appConfig.Port,
	}

	log.Fatal(srv.ListenAndServe())
}

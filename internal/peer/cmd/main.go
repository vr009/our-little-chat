package main

import (
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"our-little-chatik/internal/peer/internal/delivery"
	"our-little-chatik/internal/peer/internal/repo"
	"strconv"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type AppConfig struct {
	Port int
	DB   DBConfig
}

func main() {
	configPath := os.Getenv("PEER_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("peer-config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read a config file")
	}

	appConfig := AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     appConfig.DB.Host + ":" + appConfig.DB.Port,
		Password: appConfig.DB.Password,
	})
	peerRepo := repo.NewPeerRepository(redisClient)
	peerHandler := delivery.NewPeerHandler(peerRepo)

	diffHandler := delivery.NewDiffHandler(peerRepo, peerRepo)

	http.HandleFunc("/ws/chat", peerHandler.ConnectToChat)
	http.HandleFunc("/ws/diff", diffHandler.ConnectToDiff)

	slog.Info("service started", "port", appConfig.Port)

	err = http.ListenAndServe(":"+strconv.Itoa(appConfig.Port), nil)
	if err != nil {
		slog.Info("ListenAndServe: ", err)
	}
}

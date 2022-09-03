package main

import (
	"flag"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
	"log"
	"net/http"
	"os"
	"our-little-chatik/internal/chat_diff/internal/delivery"
	repo2 "our-little-chatik/internal/chat_diff/internal/repo"
	"our-little-chatik/internal/chat_diff/internal/usecase"
	"strconv"
)

var addr = flag.String("addr", ":8080", "http service address")

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type AppConfig struct {
	Port     int
	AuthAddr string
	DB       DBConfig
}

func main() {
	configPath := os.Getenv("CHAT_DIFF_CONFIG")
	configPath = "./internal/chat_diff/cmd"
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file")
	}

	appConfig := AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		panic(err)
	}

	ttAddr := appConfig.DB.Host + ":" + strconv.Itoa(appConfig.DB.Port)
	ttOpts := tarantool.Opts{User: appConfig.DB.Username, Pass: appConfig.DB.Password}

	conn, err := tarantool.Connect(ttAddr, ttOpts)
	if err != nil {
		panic("failed to connect to tarantool")
	}
	defer conn.Close()
	repo := repo2.NewTarantoolRepo(conn)
	chatManager := usecase.NewChatManager(repo)
	chatUsecase := usecase.NewUsecase(repo)
	tokenResolver := usecase.NewAuthResolver()
	chatServer := delivery.NewChatDiffService(chatUsecase, chatManager, tokenResolver)

	go chatManager.Work()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chatServer.WSServe(w, r)
	})

	err = http.ListenAndServe(":"+strconv.Itoa(appConfig.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

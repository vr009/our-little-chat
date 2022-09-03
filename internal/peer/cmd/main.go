package main

import (
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
	"log"
	"net/http"
	"os"
	"our-little-chatik/internal/peer/internal/delivery"
	repo2 "our-little-chatik/internal/peer/internal/repo"
	usecase2 "our-little-chatik/internal/peer/internal/usecase"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type AppConfig struct {
	Port int
	DB   DBConfig
}

func main() {
	configPath := os.Getenv("CONFIG")
	configPath = "internal/peer/cmd"
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
	messageManager := usecase2.NewMessageManager(repo)
	peerServer := delivery.NewPeerServer(messageManager)

	go messageManager.Work()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		peerServer.WSServe(w, r)
	})

	log.Printf("service started at :%d", appConfig.Port)

	err = http.ListenAndServe(":"+strconv.Itoa(appConfig.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

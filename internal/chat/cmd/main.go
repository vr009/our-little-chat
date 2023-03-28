package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"our-little-chatik/internal/chat/internal/delivery"
	"our-little-chatik/internal/chat/internal/repo"
	"our-little-chatik/internal/chat/internal/usecase"
	"our-little-chatik/internal/chat/middleware"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
)

type MongoConfig struct {
	URI      string
	Username string
	Password string
}

type AppConfig struct {
	Port int
	DB   MongoConfig
	TT   TTConfig
}

type TTConfig struct {
	Host     string
	Port     int
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

	glog.V(2)

	ttAddr := appConfig.TT.Host + ":" + strconv.Itoa(appConfig.TT.Port)
	ttOpts := tarantool.Opts{User: appConfig.DB.Username, Pass: appConfig.DB.Password}

	ttClient, err := tarantool.Connect(ttAddr, ttOpts)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ttClient.Close()

	ctx := context.Background()
	connStr, err := GetConnectionString()
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		panic(err)
	}

	repop := repo.NewPostgresRepo(conn)
	repoTT := repo.NewTarantoolRepo(ttClient)
	uc := usecase.NewChatUseCase(repop, repoTT)

	handler := delivery.NewChatHandler(uc)

	r := mux.NewRouter()

	// Getting chat messages
	r.HandleFunc("/api/v1/conv", handler.GetChatMessages).Methods("GET")
	// Getting the list of users chats
	r.HandleFunc("/api/v1/list", handler.GetChatList).Methods("GET")
	// Creating a new chat
	r.HandleFunc("/api/v1/new", handler.PostNewChat).Methods("POST")
	// Activating a chat
	r.HandleFunc("/api/v1/active", handler.PostChat).Methods("POST")

	srv := &http.Server{Handler: middleware.Logger(r),
		Addr: ":" + strconv.Itoa(appConfig.Port)}

	log.Printf("Listening port: %d", appConfig.Port)
	log.Printf("addres to query: %s", "http://localhost:"+strconv.Itoa(appConfig.Port)+"/api/v1/")
	log.Fatal(srv.ListenAndServe())
}

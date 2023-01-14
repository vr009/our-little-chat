package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"our-little-chatik/internal/chat/internal/delivery"
	"our-little-chatik/internal/chat/internal/repo"
	"our-little-chatik/internal/chat/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	ttAddr := appConfig.TT.Host + ":" + strconv.Itoa(appConfig.TT.Port)
	ttOpts := tarantool.Opts{User: appConfig.DB.Username, Pass: appConfig.DB.Password}

	ttClient, err := tarantool.Connect(ttAddr, ttOpts)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ttClient.Close()

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.DB.URI))
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database("chat_db")
	cldb := mongoClient.Database("chat_list_db")
	repom := repo.NewMongoRepo(db, cldb)
	repoTT := repo.NewTarantoolRepo(ttClient)
	uc := usecase.NewChatUseCase(repom, repoTT)

	handler := delivery.NewChatHandler(uc)

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/chat/conv", handler.GetChatMessages).Methods("GET")
	r.HandleFunc("/api/v1/chat/list", handler.GetChatList).Methods("GET")
	r.HandleFunc("/api/v1/chat/new", handler.PostNewChat).Methods("POST")
	r.HandleFunc("/api/v1/chat", handler.PostChat).Methods("POST")

	srv := &http.Server{Handler: r, Addr: ":" + strconv.Itoa(appConfig.Port)}

	log.Printf("Listening port: %d", appConfig.Port)
	log.Printf("addres to query: %s", "http://localhost:"+strconv.Itoa(appConfig.Port)+"/api/v1/")
	log.Fatal(srv.ListenAndServe())
}

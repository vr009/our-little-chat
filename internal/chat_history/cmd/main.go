package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"our-little-chatik/internal/chat_history/internal/delivery"
	"our-little-chatik/internal/chat_history/internal/repo"
	"our-little-chatik/internal/chat_history/internal/usecase"
	"strconv"
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
	configPath := os.Getenv("FLUSHER_CONFIG")
	configPath = "./internal/chat_history/cmd"
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

	ttAddr := appConfig.TT.Host + ":" + strconv.Itoa(appConfig.TT.Port)
	ttOpts := tarantool.Opts{User: appConfig.DB.Username, Pass: appConfig.DB.Password}

	ttClient, err := tarantool.Connect(ttAddr, ttOpts)
	if err != nil {
		panic("failed to connect to tarantool")
	}
	defer ttClient.Close()

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.DB.URI))
	if err != nil {
		panic(err)
	}

	db := mongoClient.Database("chat_db")
	repom := repo.NewMongoRepo(db)
	repoTT := repo.NewTarantoolRepo(ttClient)
	uc := usecase.NewChatUseCase(repom, repoTT)
	handler := delivery.NewChatHandler(uc)

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/chat/conv", handler.GetChat).Methods("GET")

	srv := &http.Server{Handler: r, Addr: ":" + strconv.Itoa(appConfig.Port)}

	log.Printf("Listening port: %d", appConfig.Port)
	log.Printf("addres to query: %s", "http://localhost:"+strconv.Itoa(appConfig.Port)+"/api/v1/")
	log.Fatal(srv.ListenAndServe())
}

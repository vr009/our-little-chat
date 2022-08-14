package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"our-little-chatik/internal/flusher/internal/delivery"
	"our-little-chatik/internal/flusher/internal/repo"
	"strconv"
)

type MongoConfig struct {
	URI      string
	Username string
	Password string
}

type TTConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type AppConfig struct {
	Port   int
	DB     MongoConfig
	TT     TTConfig
	Period int
}

// TODO 2 different mongo dbases
func main() {
	configPath := os.Getenv("FLUSHER_CONFIG")
	configPath = "./internal/flusher/cmd"
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

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	dbMsgs := mongoClient.Database("chat_db")
	dbChatList := mongoClient.Database("chat_list_db")

	msgsCol := dbMsgs.Collection("chat")
	chatListCol := dbChatList.Collection("chat_list")

	//msgsCol.Indexes().CreateOne()
	//chatListCol.Indexes().CreateOne()

	m := repo.NewMongoRepo(msgsCol, chatListCol)
	t := repo.NewTarantoolRepo(ttClient)

	daemon := delivery.NewFlusherD(t, m)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println(appConfig.Period)
	daemon.Work(ctx, appConfig.Period)
}

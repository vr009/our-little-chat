package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"our-little-chatik/internal/flusher/internal/delivery"
	"our-little-chatik/internal/flusher/internal/repo"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
)

type PostgresConfig struct {
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
	DB     PostgresConfig
	TT     TTConfig
	Period int
}

func GetConnectionString() (string, error) {
	key, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func main() {
	configPath := os.Getenv("FLUSHER_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file")
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
		log.Fatal("failed to connect to tarantool")
	}
	defer ttClient.Close()

	ctx := context.Background()
	connStr, err := GetConnectionString()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		panic(err)
	}
	m := repo.NewPostgresRepo(conn)
	t := repo.NewTarantoolRepo(ttClient)

	daemon := delivery.NewFlusherD(t, m)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println(appConfig.Period)
	daemon.Work(ctx, appConfig.Period)
}

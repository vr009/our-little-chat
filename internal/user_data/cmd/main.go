package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"our-little-chatik/internal/user_data/internal/delivery"
	"our-little-chatik/internal/user_data/internal/repo"
	"our-little-chatik/internal/user_data/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port int
}

func GetConnectionString() (string, error) {
	key, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func main() {
	configPath := os.Getenv("GATEWAY_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("user-data-config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file", configPath)
	}

	appConfig := &AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		panic(err)
	}

	connString, err := GetConnectionString()
	if err != nil {
		panic(err)
	}

	conn, err := pgxpool.Connect(context.Background(), connString)

	if err != nil {
		panic("ERROR: : " + err.Error())
	} else {
		println("Connected to postgres: ")
		println(connString)
		fmt.Println(conn.Config())
		fmt.Println(conn.Stat())
	}

	repo := repo.NewPersonRepo(conn)
	useCase := usecase.NewUserdataUseCase(repo)
	handler := delivery.NewUserdataHandler(useCase)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/user/new", handler.CreateUser).Methods("POST")

	router.HandleFunc("/api/v1/user/all", handler.GetAllUsers).Methods("GET")

	router.HandleFunc("/api/v1/user", handler.GetUser).Methods("GET")

	router.HandleFunc("/api/v1/user", handler.UpdateUser).Methods("POST")

	router.HandleFunc("/api/v1/user", handler.DeleteUser).Methods("DELETE")

	router.HandleFunc("/api/v1/user/auth", handler.CheckUserData).Methods("POST")

	srv := &http.Server{Handler: router, Addr: ":" + strconv.Itoa(appConfig.Port)}

	fmt.Printf("Main.go started at port %s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}

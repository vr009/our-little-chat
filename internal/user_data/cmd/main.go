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

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
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
	configPath := os.Getenv("USER_DATA_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("user-data-config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file", configPath)
	}

	appConfig := &AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Fatal(err)
	}

	connString, err := GetConnectionString()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.Connect(context.Background(), connString)

	if err != nil {
		log.Fatal("ERROR: : " + err.Error())
	} else {
		glog.Infof("Connected to postgres: %s", connString)
	}

	personRepo := repo.NewPersonRepo(conn)
	useCase := usecase.NewUserdataUseCase(personRepo)
	userDataHandler := delivery.NewUserdataHandler(useCase)
	authHandler := delivery.NewAuthHandler(useCase)

	router := mux.NewRouter()

	// Administration API
	router.HandleFunc("/api/v1/user/new", userDataHandler.CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/user/all", userDataHandler.GetAllUsers).Methods("GET")
	router.HandleFunc("/api/v1/user", userDataHandler.UpdateUser).Methods("POST")
	router.HandleFunc("/api/v1/user", userDataHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/v1/user/auth", userDataHandler.CheckUserData).Methods("POST")

	// Common API
	router.HandleFunc("/api/v1/user/me", userDataHandler.GetUser).Methods("GET")
	router.HandleFunc("/api/v1/user/search", userDataHandler.FindUser).Methods("GET")

	// Auth API
	router.HandleFunc("/api/v1/auth/signup", authHandler.SignUp).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/v1/auth/logout", authHandler.Logout).Methods("DELETE")

	srv := &http.Server{Handler: router, Addr: ":" + strconv.Itoa(appConfig.Port)}

	fmt.Printf("Main.go started at port %s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}

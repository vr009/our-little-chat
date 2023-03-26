package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"our-little-chatik/internal/auth/internal/delivery"
	"our-little-chatik/internal/auth/internal/models"
	"our-little-chatik/internal/auth/internal/repo"
	"our-little-chatik/internal/auth/internal/usecase"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var defaultLevelLog glog.Level = 2

func main() {
	configPath := os.Getenv("AUTH_CONFIG")

	viper.AddConfigPath(configPath)
	viper.SetConfigName("auth-config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found; ignore error if desired")
		} else {
			fmt.Println("Config file was found but another error was produced")
		}
	}

	glog.V(defaultLevelLog)

	appConfig := models.AppConfig{}

	if err := viper.Unmarshal(&appConfig); err != nil {
		fmt.Println(err)
		log.Fatal("Error of unmarshal")
	}

	dbInfo := redis.Options{
		Addr:     appConfig.DataBase.Port,
		Password: appConfig.DataBase.Password,
		DB:       appConfig.DataBase.DB,
	}

	client := redis.NewClient(&dbInfo)

	if client == nil {
		log.Fatal("client doesnt work")
	}

	glog.Infof("Redis started at port %s \n", appConfig.DataBase.Port)

	db := repo.NewDataBase(client, appConfig.DataBase.TtlHours)
	useCase := usecase.NewAuthUseCase(db)
	handler := delivery.NewAuthHandler(useCase)

	router := mux.NewRouter()

	// Получение Token пользователя по UserID
	// (UserID) => Token
	router.HandleFunc("/api/v1/auth/token", handler.GetToken).Methods("GET")

	// Получение UserID по Token
	// (Token) => UserID
	router.HandleFunc("/api/v1/auth/user", handler.GetUser).Methods("GET")

	// Добавление нового UserID и создание Token
	// (UserID) => Session {
	//	   UserID: Token
	//	   Token: UserID
	//	}
	router.HandleFunc("/api/v1/auth", handler.PostSession).Methods("POST")

	// 4.
	// Удаление сессии по Token
	// Token –> Session {}
	router.HandleFunc("/api/v1/auth", handler.DeleteSession).Methods("DELETE")

	srv := &http.Server{Handler: router, Addr: appConfig.Address}

	glog.Infof("Main.go started at port %s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}

package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"our-little-chatik/internal/gateway/internal/delivery"
	"our-little-chatik/internal/gateway/internal/usecase"
	"our-little-chatik/internal/models"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port int
}

func main() {
	configPath := os.Getenv("GATEWAY_CONFIG")
	configPath = "./internal/gateway/cmd"
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file")
	}

	appConfig := &AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	srv := &http.Server{
		Handler: r,
		Addr:    strconv.Itoa(appConfig.Port),
	}

	userDataCfg := models.ServiceRouterConfig{
		BaseUrl: "http://localhost:8086/api/v1",
		Router: map[string]string{
			"AddUser":    "/api/v1/user/new",
			"GetUser":    "/api/v1/user",
			"DeleteUser": "/api/v1/user",
		},
	}
	userDataClient := http.Client{}
	userDataHandler := delivery.NewUserDataHandler(userDataClient, userDataCfg)

	authCfg := models.ServiceRouterConfig{
		BaseUrl: "http://localhost:8086/api/v1",
		Router: map[string]string{
			"AddUser":    "/api/v1/auth",
			"GetUser":    "/api/v1/auth/token",
			"DeleteUser": "/api/v1/auth",
			"GetSession": "/api/v1/auth/user",
		},
	}
	authClient := http.Client{}
	authHandler := delivery.NewAuthHandler(authClient, authCfg)

	uc := usecase.NewGatewayUsecasse(*authHandler, *userDataHandler)
	gatewayHandler := delivery.NewGatewayHandler(uc)

	r.HandleFunc("/api/gateway/signup", gatewayHandler.SignUp).Methods("POST")
	r.HandleFunc("/api/gateway/signin", gatewayHandler.SignIn).Methods("POST")
	r.HandleFunc("/api/gateway/session", gatewayHandler.GetSession).Methods("GET")
	r.HandleFunc("/api/gateway/user", gatewayHandler.GetUser).Methods("GET")

	log.Fatal(srv.ListenAndServe())
}

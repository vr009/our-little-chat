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
	Port            int
	UserDataBaseURL string
	AuthBaseURL     string
}

func main() {
	configPath := os.Getenv("GATEWAY_CONFIG")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("gateway-config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read a config file", configPath)
	}

	appConfig := &AppConfig{}
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + strconv.Itoa(appConfig.Port),
	}

	userDataCfg := models.ServiceRouterConfig{
		BaseUrl: appConfig.UserDataBaseURL,
		Router: map[string]string{
			"AddUser":    "/user/new",
			"CheckUser":  "/user/auth",
			"DeleteUser": "/user",
			"FindUser":   "/user/search",
		},
	}
	userDataClient := http.Client{}
	userDataHandler := delivery.NewUserDataHandler(userDataClient, &userDataCfg)

	authCfg := models.ServiceRouterConfig{
		BaseUrl: appConfig.AuthBaseURL,
		Router: map[string]string{
			"AddUser":    "/auth",
			"GetUser":    "/auth/user",
			"DeleteUser": "/auth",
			"GetSession": "/auth/token",
		},
	}
	authClient := http.Client{}
	authHandler := delivery.NewAuthHandler(authClient, &authCfg)

	uc := usecase.NewGatewayUsecasse(*authHandler, *userDataHandler)
	gatewayHandler := delivery.NewGatewayHandler(uc)

	r.HandleFunc("/api/gateway/signup", gatewayHandler.SignUp).Methods("POST")
	r.HandleFunc("/api/gateway/signin", gatewayHandler.SignIn).Methods("POST")
	r.HandleFunc("/api/gateway/logout", gatewayHandler.LogOut).Methods("DELETE")
	r.HandleFunc("/api/gateway/search", gatewayHandler.Find).Methods("GET")

	log.Fatal(srv.ListenAndServe())
}

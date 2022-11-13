package main

import (
	"log"
	"net/http"
	"os"

	"our-little-chatik/internal/gateway/internal/delivery"
	"our-little-chatik/internal/gateway/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port           int
	UserDataConfig delivery.UserDataConfig
	AuthConfig     delivery.AuthServiceConfig
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
		Addr:    ":8080",
	}

	userDataCfg := delivery.UserDataConfig{
		BaseUrl: "",
		Router:  map[string]string{},
	}
	userDataClient := http.Client{}
	userDataHandler := delivery.NewUserDataHandler(userDataClient, userDataCfg)

	authCfg := delivery.AuthServiceConfig{
		BaseUrl: "",
		Router:  map[string]string{},
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

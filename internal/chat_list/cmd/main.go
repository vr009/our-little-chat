package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/gorilla/docs"
	"github.com/tarantool/go-tarantool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"our-little-chatik/internal/chat_list/internal/delivery"
	repo2 "our-little-chatik/internal/chat_list/internal/repo"
	usecase2 "our-little-chatik/internal/chat_list/internal/usecase"
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

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /
func main() {
	configPath := os.Getenv("CHAT_LIST_CONFIG")
	configPath = "./internal/chat_list/cmd"
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

	db := mongoClient.Database("chat_list_db")
	col := db.Collection("chat_list")

	repo := repo2.NewChatListRepo(col)
	usecase := usecase2.NewChatListUsecase(repo)
	handler := delivery.NewChatListHandler(usecase)

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/chats", handler.GetChatList).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	srv := &http.Server{Handler: r, Addr: ":" + strconv.Itoa(appConfig.Port)}

	log.Printf("listening port: %d", appConfig.Port)
	log.Printf("addres to query: %s", "http://localhost:"+strconv.Itoa(appConfig.Port)+"/api/v1/")
	log.Fatal(srv.ListenAndServe())
}

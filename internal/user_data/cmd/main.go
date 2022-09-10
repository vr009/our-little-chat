package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"our-little-chatik/internal/user_data/internal/delivery"
	"our-little-chatik/internal/user_data/internal/repo"
	"our-little-chatik/internal/user_data/internal/usecase"
)

func init() {
	godotenv.Load(".env")
}

func GetConnectionString() (string, error) {
	key, flag := os.LookupEnv("DATABASE_URL")
	if !flag {
		return "", errors.New("connection string not found")
	}
	return key, nil
}

func main() {

	fmt.Println("Starting..")

	connString := "user=postgres password=postgres host=localhost port=5432 dbname=postgres"

	conn, err := pgxpool.Connect(context.Background(), connString)

	if err != nil {
		println("ERROR: : " + err.Error())
	} else {
		println("Connected to postgres: ")
		println(connString)
	}

	repo := repo.NewPersonRepo(conn)

	useCase := usecase.NewUserdataUseCase(repo)

	handler := delivery.NewUserdataHandler(useCase)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/user_data/create", handler.CreateUser).Methods("POST")

	router.HandleFunc("/api/v1/user_data/all_users", handler.GetAllUsers).Methods("GET")

	router.HandleFunc("/api/v1/user_data/user", handler.GetUser).Methods("GET")

	router.HandleFunc("/api/v1/user_data/user", handler.UpdateUser).Methods("POST")

	router.HandleFunc("/api/v1/user_data/user", handler.DeleteUser).Methods("DELETE")

	srv := &http.Server{Handler: router, Addr: ":8080"}

	fmt.Printf("Main.go started at port %s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}

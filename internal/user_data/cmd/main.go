package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"our-little-chatik/internal/user_data/internal/delivery"
	"our-little-chatik/internal/user_data/internal/repo"
	"our-little-chatik/internal/user_data/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
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

	connString, err := GetConnectionString()
	if err != nil {
		panic("failed to get a connection string")
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

	srv := &http.Server{Handler: router, Addr: ":8086"}

	fmt.Printf("Main.go started at port %s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}

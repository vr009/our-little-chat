package main

import (
	"auth/internal/delivery"
	repo2 "auth/internal/repo"
	"auth/internal/usecase"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Starting..")

	dbInfo := redis.Options{
		Addr:     ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	client := redis.NewClient(&dbInfo)

	if client == nil {
		panic("client doesnt work")
	}

	fmt.Printf("Redis started at port %s \n", dbInfo.Addr)

	repo := repo2.NewDataBase(client)
	useCase := usecase.NewAuthUseCase(repo)
	handler := delivery.NewAuthHandler(useCase)

	router := mux.NewRouter()

	// 1. – РАБОТАЕТ!)
	// Получение Token пользователя по UserID
	// (UserID) => Token
	router.HandleFunc("/api/v1/auth/token", handler.GetToken).Methods("GET")

	// 2. – РАБОТАЕТ!)
	// Получение UserID по Token
	// (Token) => UserID
	router.HandleFunc("/api/v1/auth/user", handler.GetUser).Methods("GET")

	// 3. – РАБОТАЕТ!)
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

	srv := &http.Server{Handler: router, Addr: ":8080"}

	fmt.Printf("Main.go started at port %s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}

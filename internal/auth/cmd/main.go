package main

import (
	"auth/internal/delivery"
	repo2 "auth/internal/repo"
	"auth/internal/usecase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	repo := repo2.NewDataBase()
	useCase := usecase.NewAuthUseCase(repo)
	handler := delivery.NewAuthHandler(useCase)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/auth", handler.GetSession).Methods("GET")
	router.HandleFunc("/api/v1/auth", handler.PostSession).Methods("POST")
	router.HandleFunc("/api/v1/auth", handler.DeleteSession).Methods("DELETE")

	srv := &http.Server{Handler: router, Addr: ":8080"}

	log.Fatal(srv.ListenAndServe())
}

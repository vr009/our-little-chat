package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

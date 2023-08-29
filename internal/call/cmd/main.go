package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
	"our-little-chatik/internal/call/internal"
)

func main() {
	internal.AllRooms.Init()

	port := os.Getenv("CALL_PORT")
	if port == "" {
		panic("empty CALL_PORT provided")
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	//http.HandleFunc("/create", internal.CreateRoomRequestHandler)
	//http.HandleFunc("/join", internal.JoinRoomRequestHandler)

	r := mux.NewRouter()
	r.HandleFunc("/call", internal.CreateOrJoinRoomHandler)
	r.HandleFunc("/finish", internal.DeleteRoomHandler).Methods("DELETE")

	slog.Info(fmt.Sprintf("Starting Server on Port %s", port))
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}

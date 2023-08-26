package main

import (
	"fmt"
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

	http.HandleFunc("/create", internal.CreateRoomRequestHandler)
	http.HandleFunc("/join", internal.JoinRoomRequestHandler)

	slog.Info(fmt.Sprintf("Starting Server on Port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

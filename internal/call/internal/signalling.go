package internal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// AllRooms is the global hashmap for the server
var AllRooms RoomMap

// CreateRoomRequestHandler Create a Room and return roomID
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	chatID := r.URL.Query().Get("chat_id")
	roomID := AllRooms.CreateRoom(chatID)

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				log.Println("!!!!!!", client.Conn, msg.Client)
				client.Mutex.Lock()
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					client.Mutex.Unlock()
					client.Conn.Close()
					log.Println(err)
					return
				}
				client.Mutex.Unlock()
			}
		}
	}
}

// JoinRoomRequestHandler will join the client in a particular room
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomID")

	if roomID == "" {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomID, false, ws)

	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Println("Read Error: ", err)
			return
		}

		msg.Client = ws
		msg.RoomID = roomID

		log.Println(msg.Message)

		broadcast <- msg
	}
}

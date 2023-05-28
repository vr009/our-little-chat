package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"our-little-chatik/internal/chat_diff/internal"
	models2 "our-little-chatik/internal/chat_diff/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type ChatDiffService struct {
	manager internal.Manager
}

func NewChatDiffService(manager internal.Manager) *ChatDiffService {
	return &ChatDiffService{manager: manager}
}

type WebSocketClient struct {
	conn        *websocket.Conn
	currentUser *models2.ChatDiffUser
	manager     internal.Manager

	// Buffered channel of outbound messages.
	send chan []byte
}

func newWebSocketClient(conn *websocket.Conn, manager internal.Manager) *WebSocketClient {
	client := &WebSocketClient{conn: conn, manager: manager}
	return client
}

func (ws *WebSocketClient) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.conn.Close()
	}()
	for {
		if ws.currentUser == nil {
			time.Sleep(1)
			continue
		}
		select {
		case chatUpdates := <-ws.currentUser.Updates:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			for _, upd := range chatUpdates {
				buf, err := json.Marshal(upd)
				if err != nil {
					log.Fatal(err)
					return
				}
				w.Write(buf)
				w.Write(newline)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// read pumps messages from the websocket connection.
//
// The application runs read in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (ws *WebSocketClient) read() {
	defer func() {
		ws.conn.Close()
	}()
	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		if ws.currentUser != nil {
			time.Sleep(1)
			continue
		}
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		user := &models2.User{}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		err = json.Unmarshal(message, user)
		if err != nil {
			log.Fatalf("failed to unmarshal message")
		}

		chatUser := &models2.ChatDiffUser{User: *user}

		ws.currentUser = ws.manager.AddChatUser(chatUser)

		fmt.Println("returned", ws.currentUser, &ws.currentUser)
		log.Println(user)
	}
}

func (server *ChatDiffService) WSServe(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := newWebSocketClient(conn, server.manager)
	go client.write()
	go client.read()
}

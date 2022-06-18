package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"our-little-chatik/internal/peer/internal"
	"our-little-chatik/internal/peer/internal/models"
	"time"
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

type PeerServer struct {
	uc      internal.PeerUsecase
	manager internal.MessageManager
}

func NewPeerServer(uc internal.PeerUsecase, manager internal.MessageManager) *PeerServer {
	return &PeerServer{
		uc:      uc,
		manager: manager,
	}
}

type WebSocketClient struct {
	conn        *websocket.Conn
	uc          internal.PeerUsecase
	currentChat *models.Chat
	manager     internal.MessageManager

	// Buffered channel of outbound messages.
	send chan []byte
}

func newWebSocketClient(conn *websocket.Conn, uc internal.PeerUsecase, manager internal.MessageManager) *WebSocketClient {
	client := &WebSocketClient{conn: conn, uc: uc, manager: manager}
	return client
}

func (ws *WebSocketClient) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.conn.Close()
	}()
	for {
		if ws.currentChat == nil {
			time.Sleep(1)
			continue
		}
		select {
		case messages := <-ws.currentChat.ReadyForRecv:
			//log.Println("sending for receiver:", ws.currentChat.ReceiverID)
			//log.Println(messages)
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			for _, msg := range messages {
				buf, err := json.Marshal(msg)
				if err != nil {
					log.Fatal(err)
					return
				}
				w.Write(buf)
				w.Write(newline)
			}

			//log.Println("sent.")

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
		log.Println("closed connection")
		ws.conn.Close()
		if ws.currentChat != nil && ws.currentChat.ChatSibling != nil {
			ws.manager.DequeueChat(ws.currentChat.ChatSibling)
		}
		if ws.currentChat != nil {
			ws.manager.DequeueChat(ws.currentChat)
		}
	}()

	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg := &models.Message{}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		err = json.Unmarshal(message, msg)
		if err != nil {
			log.Println("failed to unmarshal message")
			return
		}

		chat := models.NewChatFromMsg(msg)
		chat = ws.manager.EnqueueChatIfNotExists(chat)
		if ws.currentChat == nil {
			ws.currentChat = chat
		}

		msgForSibling := models.NewMessageForAnotherSide(msg)
		if chat.ChatSibling == nil {
			chat.ChatSibling = models.NewChatFromMsg(msgForSibling)
			chat.ChatSibling = ws.manager.EnqueueChatIfNotExists(chat.ChatSibling)
		}
		chat.ChatSibling.ChatSibling = chat

		if msg.SessionStart {
			continue
		}

		err = ws.uc.SendMessage(msg, chat)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (server *PeerServer) WSServe(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := newWebSocketClient(conn, server.uc, server.manager)
	go client.write()
	go client.read()
}

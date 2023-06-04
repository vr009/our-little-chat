package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"our-little-chatik/internal/peer/internal"
	"our-little-chatik/internal/peer/internal/models"

	"github.com/golang/glog"
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

type PeerServer struct {
	manager internal.MessageManager
}

func NewPeerServer(manager internal.MessageManager) *PeerServer {
	return &PeerServer{
		manager: manager,
	}
}

type WebSocketClient struct {
	conn        *websocket.Conn
	currentPeer *models.Peer
	currentChat *models.Chat
	manager     internal.MessageManager

	// Buffered channel of outbound messages.
	send chan []byte

	disconnected chan struct{}
}

func newWebSocketClient(conn *websocket.Conn, manager internal.MessageManager) *WebSocketClient {
	client := &WebSocketClient{conn: conn, manager: manager}
	return client
}

func (ws *WebSocketClient) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.conn.Close()
		fmt.Println("Write closed")
	}()
	for {
		if ws.currentPeer == nil {
			time.Sleep(1)
			continue
		}
		select {
		case messages := <-ws.currentPeer.MsgsToRecv:
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
		case <-ws.disconnected:
			return
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
		if ws.currentPeer != nil {
			err := ws.currentChat.UnsubscribePeer(ws.currentPeer)
			log.Println(err)
		}
		ws.disconnected <- struct{}{}
		log.Println("closed connection")
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
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("failed to unmarshal message", err)
			continue
		}
		fmt.Println("Received", msg)
		if ws.currentPeer != nil && !msg.SessionStart {
			ws.currentPeer.MsgToSend <- msg
			continue
		}

		if msg.SessionStart {
			newChat := models.GetChatFromMessage(msg)
			chat := ws.manager.EnqueueChatIfNotExists(newChat)
			peer := models.GetPeerFromMessage(msg)
			err = chat.SubscribePeer(peer)
			if err != nil {
				return
			}
			ws.currentChat = chat
			ws.currentPeer = peer
			fmt.Println("subscribed peer: ", peer.PeerID)
		}
	}
}

func (server *PeerServer) WSServe(w http.ResponseWriter, r *http.Request) {
	//user, err := pkg.AuthHook(r)
	//if err != nil {
	//	w.WriteHeader(http.StatusForbidden)
	//	errObj := models2.Error{Msg: "Invalid token"}
	//	body, _ := json.Marshal(errObj)
	//	w.Write(body)
	//	glog.Error(errObj.Msg)
	//	return
	//}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Error(err.Error())
		return
	}
	client := newWebSocketClient(conn, server.manager)
	//glog.Infof("Created new connection:\nID: %s", user.UserID)
	client.disconnected = make(chan struct{})
	go client.write()
	go client.read()
}

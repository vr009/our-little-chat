package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"log"
	"net/http"
	"our-little-chatik/internal/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"our-little-chatik/internal/chat_diff/internal"
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

type PeerHandler struct {
	peersMap sync.Map
	repo     internal.DiffRepo
}

func NewPeerHandler(repo internal.DiffRepo) *PeerHandler {
	return &PeerHandler{
		repo: repo,
	}
}

func (h *PeerHandler) ConnectToChats(w http.ResponseWriter, r *http.Request) {
	//user := strings.TrimPrefix(r.URL.Path, "/chat/")

	chat_id := r.Form.Get("chat_id")
	user_id := r.Form.Get("user_id")

	if chat_id == "" || user_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peer, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("websocket conn failed", err)
	}

}

// ChatSession represents a connected/active chat user
type ChatSession struct {
	user     string
	peer     *websocket.Conn
	diffRepo internal.DiffRepo
}

// NewChatSession returns a new ChatSession
func NewChatSession(user string, peer *websocket.Conn,
	repo internal.DiffRepo) *ChatSession {
	return &ChatSession{user: user, peer: peer, repo: repo}
}

const userSet = "%s_%s"

// Start starts the chat by reading messages sent by the peer and broadcasting the to redis pub-sub channel
func (s *ChatSession) Start() {
	usernameTaken, err := s.repo.CheckUserExists(context.Background(),
		s.user, fmt.Sprintf(userSet, "users", s.chatID))

	if err != nil {
		log.Println("unable to determine whether user exists -", s.user)
		s.notifyPeer(retryMessage)
		s.peer.Close()
		return
	}

	if usernameTaken {
		msg := fmt.Sprintf(usernameHasBeenTaken, s.user)
		s.peer.WriteMessage(websocket.TextMessage, []byte(msg))
		s.peer.Close()
		return
	}

	err = s.repo.CreateUser(context.Background(),
		s.user, fmt.Sprintf(userSet, "users", s.chatID))
	if err != nil {
		log.Println("failed to add user to list of active chat users", s.user)
		s.notifyPeer(retryMessage)
		s.peer.Close()
		return
	}
	s.Peers[s.user] = s.peer

	s.notifyPeer(fmt.Sprintf(welcome, s.user))

	/*
		this go-routine will exit when:
		(1) the user disconnects from chat manually
		(2) the app is closed
	*/
	go func() {
		log.Println("user joined", s.user)
		for {
			_, bMsg, err := s.peer.ReadMessage()
			if err != nil {
				_, ok := err.(*websocket.CloseError)
				if ok {
					log.Println("connection closed by user")
					s.disconnect()
				}
				return
			}

			chatID, err := uuid.Parse(s.chatID)
			if err != nil {
				glog.Error(err)
				return
			}
			senderID, err := uuid.Parse(s.user)
			if err != nil {
				glog.Error(err)
				return
			}

			msg := models.Message{
				MsgID:     uuid.New(),
				Payload:   string(bMsg),
				ChatID:    chatID,
				SenderID:  senderID,
				CreatedAt: time.Now().Unix(),
			}
			s.repo.SendToChannel(context.Background(),
				msg, fmt.Sprintf(userSet, "users", s.chatID))
			// persist message
			err = s.repo.SaveMessage(msg)
			if err != nil {
				glog.Error(err)
			}
		}
	}()
	go func() {
		msgChan := make(chan models.Message)
		s.repo.StartSubscriber(context.Background(),
			msgChan, s.chatID)
		for {
			select {
			case msg := <-msgChan:
				fmt.Printf("got your message: %s from %s\n", msg.Payload, msg.SenderID.String())
				for user, peer := range s.Peers {
					if msg.SenderID.String() != user { //don't recieve your own messages
						bMsg, err := json.Marshal(msg)
						if err != nil {
							glog.Error(err)
							break
						}
						peer.WriteMessage(websocket.TextMessage, bMsg)
					}
				}
			}
		}
	}()
}

func (s *ChatSession) notifyPeer(msg string) {
	err := s.peer.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("failed to write message", err)
	}
}

// Invoked when the user disconnects (websocket connection is closed). It performs cleanup activities
func (s *ChatSession) disconnect() {
	//remove user from SET
	s.repo.RemoveUser(context.Background(),
		s.user, fmt.Sprintf(userSet, "users", s.chatID))

	//close websocket
	s.peer.Close()

}

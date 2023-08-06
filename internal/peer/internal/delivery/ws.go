package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/peer/internal"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type PeerHandler struct {
	peersMap sync.Map
	repo     internal.PeerRepo
}

func NewPeerHandler(repo internal.PeerRepo) *PeerHandler {
	return &PeerHandler{
		repo: repo,
	}
}

func (h *PeerHandler) ConnectToChat(w http.ResponseWriter, r *http.Request) {
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

	var peers map[string]*websocket.Conn
	var val interface{}
	var ok bool

	if val, ok = h.peersMap.Load(chat_id); !ok {
		peers = make(map[string]*websocket.Conn)
	} else {
		peers = val.(map[string]*websocket.Conn)
		h.peersMap.Store(chat_id, peers)
	}

	chatSession := NewChatSession(user_id, peer, peers, chat_id, h.repo)
	chatSession.Start()
}

// ChatSession represents a connected/active chat user
type ChatSession struct {
	user   string
	peer   *websocket.Conn
	Peers  map[string]*websocket.Conn
	repo   internal.PeerRepo
	chatID string
}

// NewChatSession returns a new ChatSession
func NewChatSession(user string, peer *websocket.Conn,
	peers map[string]*websocket.Conn, chatID string,
	repo internal.PeerRepo) *ChatSession {
	return &ChatSession{user: user, peer: peer, Peers: peers, repo: repo}
}

const usernameHasBeenTaken = "username %s is already taken. please retry with a different name"
const retryMessage = "failed to connect. please try again"
const welcome = "Welcome %s!"

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

	//remove from Peers
	delete(s.Peers, s.user)
}

package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/peer/internal"
	models2 "our-little-chatik/internal/peer/internal/models"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type PeerHandler struct {
	// peersMap is a sync map where keys are chat ids and values
	// are maps where keys are user ids and values are their corresponding
	// websocket connections
	peersMap sync.Map
	repo     internal.PeerRepo
	msgBus   internal.MessageBus
}

func NewPeerHandler(repo internal.PeerRepo, msgBus internal.MessageBus) *PeerHandler {
	return &PeerHandler{
		repo:   repo,
		msgBus: msgBus,
	}
}

func (h *PeerHandler) ConnectToChat(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")
	userID := r.URL.Query().Get("user_id")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	if chatID == "" || userID == "" {
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

	if val, ok = h.peersMap.Load(chatID); !ok {
		peers = make(map[string]*websocket.Conn)
		h.peersMap.Store(chatID, peers)
	} else {
		peers = val.(map[string]*websocket.Conn)
	}

	chatSession := NewChatSession(userID, peer, peers, chatID, h.repo, h.msgBus)
	chatSession.Start()
}

// ChatSession represents a connected/active chat user
type ChatSession struct {
	user   string
	peer   *websocket.Conn
	Peers  map[string]*websocket.Conn
	repo   internal.PeerRepo
	chatID string
	msgBus internal.MessageBus
}

// NewChatSession returns a new ChatSession
func NewChatSession(user string, peer *websocket.Conn,
	peers map[string]*websocket.Conn, chatID string,
	repo internal.PeerRepo, msgBus internal.MessageBus) *ChatSession {
	return &ChatSession{
		user:   user,
		peer:   peer,
		Peers:  peers,
		repo:   repo,
		chatID: chatID,
		msgBus: msgBus,
	}
}

const usernameHasBeenTaken = "username %s is already taken. please retry with a different name"
const retryMessage = "failed to connect. please try again"
const welcome = "Welcome %s!"

// Start starts the chat by reading messages sent by the peer and broadcasting the to redis pub-sub channel
func (s *ChatSession) Start() {
	usernameTaken, err := s.repo.CheckUserExists(context.Background(),
		s.user, fmt.Sprintf(models2.ChatUsersFmtStr, "users", s.chatID))
	if err != nil {
		log.Println("unable to determine whether user exists -", s.user)
		s.notifyPeer(retryMessage)
		s.peer.Close()
		return
	}
	if usernameTaken {
		msg := fmt.Sprintf(usernameHasBeenTaken, s.user)
		s.notifyPeer(msg)
		s.peer.Close()
		return
	}

	err = s.repo.CreateUser(context.Background(),
		s.user, fmt.Sprintf(models2.ChatUsersFmtStr, "users", s.chatID))
	if err != nil {
		log.Println("failed to add user to list of active chat users", s.user)
		s.notifyPeer(retryMessage)
		s.peer.Close()
		return
	}
	s.Peers[s.user] = s.peer

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
				log.Println("=========== disconnecting ===========")
				_, ok := err.(*websocket.CloseError)
				if ok {
					log.Println("connection closed by user")
					s.disconnect()
				}
				return
			}

			chatID, err := uuid.Parse(s.chatID)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			senderID, err := uuid.Parse(s.user)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			msg := models.Message{
				MsgID:     uuid.New(),
				Payload:   string(bMsg),
				ChatID:    chatID,
				SenderID:  senderID,
				CreatedAt: time.Now().Unix(),
			}
			// persist message
			err = s.repo.SaveMessage(msg)
			if err != nil {
				slog.Error(err.Error())
			}
			// Send via message bus
			s.msgBus.SendMessageToChannel(context.Background(),
				msg, fmt.Sprintf(models2.ChatUsersFmtStr, "users", s.chatID))
		}
	}()

	readyChan := make(chan struct{})
	go func() {
		msgChan := s.msgBus.SubscribeOnChatMessages(context.Background(),
			fmt.Sprintf(models2.ChatUsersFmtStr, "users", s.chatID), readyChan)

		for {
			select {
			case msg := <-msgChan:
				// TODO probably, there is no need to send it to all peers and hold the map at all.
				//  It is enough to send to the current peer - other peers will receive their
				//  messages in another goroutine
				for _, peer := range s.Peers {
					err := sendMessageToPeer(peer, msg)
					if err != nil {
						slog.Error(err.Error())
						break
					}
				}
			}
		}
	}()

	<-readyChan
	s.notifyPeer(fmt.Sprintf(welcome, s.user))
}

func sendMessageToPeer(peer *websocket.Conn, msg models.Message) error {
	notification := models2.Notification{
		Type:    models2.ChatMessage,
		Message: &msg,
	}
	bMsg, err := json.Marshal(&notification)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	peer.WriteMessage(websocket.TextMessage, bMsg)
	return nil
}

func (s *ChatSession) notifyPeer(msg string) {
	notification := models2.Notification{
		Type:        models2.InfoMessage,
		Description: msg,
	}
	bNotification, _ := json.Marshal(notification)
	err := s.peer.WriteMessage(websocket.TextMessage, bNotification)
	if err != nil {
		log.Println("failed to write message", err)
	}
}

// Invoked when the user disconnects (websocket connection is closed). It performs cleanup activities
func (s *ChatSession) disconnect() {
	//remove user from SET
	s.repo.RemoveUser(context.Background(),
		s.user, fmt.Sprintf(models2.ChatUsersFmtStr, "users", s.chatID))

	//close websocket
	s.peer.Close()

	//remove from Peers
	delete(s.Peers, s.user)
}

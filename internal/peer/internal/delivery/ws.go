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
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type PeerHandler struct {
	repo   internal.PeerRepo
	msgBus internal.MessageBus
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

	chatSession := NewChatSession(userID, peer, chatID, h.repo, h.msgBus)
	chatSession.Start()
}

// ChatSession represents a connected/active chat user
type ChatSession struct {
	userID   string
	peerConn *websocket.Conn
	repo     internal.PeerRepo
	chatID   string
	msgBus   internal.MessageBus
}

// NewChatSession returns a new ChatSession
func NewChatSession(userID string, peerConn *websocket.Conn, chatID string,
	repo internal.PeerRepo, msgBus internal.MessageBus) *ChatSession {
	return &ChatSession{
		userID:   userID,
		peerConn: peerConn,
		chatID:   chatID,
		repo:     repo,
		msgBus:   msgBus,
	}
}

const usernameHasBeenTaken = "username %s is already taken. please retry with a different name"
const retryMessage = "failed to connect. please try again"
const welcome = "Welcome %s!"

// Start starts the chat by reading messages sent by the peer and broadcasting the to redis pub-sub channel
func (s *ChatSession) Start() {
	usernameTaken, err := s.repo.CheckUserExists(context.Background(),
		s.userID, fmt.Sprintf(models2.CommonFormat, "users", s.userID))
	if err != nil {
		log.Println("unable to determine whether user exists -", s.userID)
		s.notifyPeer(retryMessage)
		s.peerConn.Close()
		return
	}
	if usernameTaken {
		msg := fmt.Sprintf(usernameHasBeenTaken, s.userID)
		s.notifyPeer(msg)
		s.peerConn.Close()
		return
	}

	err = s.repo.CreateUser(context.Background(),
		s.userID, fmt.Sprintf(models2.CommonFormat, "users", s.userID))
	if err != nil {
		log.Println("failed to add user to list of active chat users", s.userID)
		s.notifyPeer(retryMessage)
		s.peerConn.Close()
		return
	}

	/*
		this go-routine will exit when:
		(1) the user disconnects from chat manually
		(2) the app is closed
	*/
	go func() {
		log.Println("user joined", s.userID)
		for {
			_, bMsg, err := s.peerConn.ReadMessage()
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
			senderID, err := uuid.Parse(s.userID)
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
				msg, fmt.Sprintf(models2.CommonFormat, "chat", s.chatID))
		}
	}()

	readyChan := make(chan struct{})
	go func() {
		msgChan := s.msgBus.SubscribeOnChatMessages(context.Background(),
			fmt.Sprintf(models2.CommonFormat, "chat", s.chatID), readyChan)

		for {
			select {
			case msg := <-msgChan:
				err := sendMessageToPeer(s.peerConn, msg)
				if err != nil {
					slog.Error(err.Error())
					break
				}
			}
		}
	}()

	<-readyChan
	s.notifyPeer(fmt.Sprintf(welcome, s.userID))
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
	err := s.peerConn.WriteMessage(websocket.TextMessage, bNotification)
	if err != nil {
		log.Println("failed to write message", err)
	}
}

// Invoked when the user disconnects (websocket connection is closed). It performs cleanup activities
func (s *ChatSession) disconnect() {
	//remove user from SET
	s.repo.RemoveUser(context.Background(),
		s.userID, fmt.Sprintf(models2.CommonFormat, "users", s.userID))

	//close websocket
	s.peerConn.Close()
}

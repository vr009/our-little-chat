package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/peer/internal"
	models2 "our-little-chatik/internal/peer/internal/models"
	"sync"
)

type DiffHandler struct {
	peersMap sync.Map
	repo     internal.PeerRepo
	diffRepo internal.DiffRepo
}

func NewDiffHandler(repo internal.PeerRepo, diffRepo internal.DiffRepo) *DiffHandler {
	return &DiffHandler{
		repo:     repo,
		diffRepo: diffRepo,
	}
}

func (h *DiffHandler) ConnectToDiff(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	peer, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("websocket conn failed", err)
	}

	chatSession := NewDiffSession(userID, peer, h.repo, h.diffRepo)
	chatSession.Start()
}

// DiffSession represents a connected/active chats diff user
type DiffSession struct {
	userID   string
	peer     *websocket.Conn
	repo     internal.PeerRepo
	diffRepo internal.DiffRepo
}

// NewDiffSession returns a new DiffSession
func NewDiffSession(userID string, peer *websocket.Conn,
	repo internal.PeerRepo, diffRepo internal.DiffRepo) *DiffSession {
	return &DiffSession{userID: userID, peer: peer, repo: repo, diffRepo: diffRepo}
}

// Start starts the chat by reading messages sent by the peer and broadcasting the to redis pub-sub channel
func (s *DiffSession) Start() {
	usernameTaken, err := s.repo.CheckUserExists(context.Background(),
		s.userID, fmt.Sprintf(models2.CommonFormat, "diff_users", s.userID))
	if err != nil {
		log.Println("unable to determine whether user exists -", s.userID)
		s.notifyPeer(models2.Failed, map[string]any{
			"description": retryMessage,
		})
		s.peer.Close()
		return
	}
	if usernameTaken {
		msg := fmt.Sprintf(usernameHasBeenTaken, s.userID)
		s.notifyPeer(models2.Conflict, map[string]any{
			"description": msg,
		})
		s.peer.Close()
		return
	}

	err = s.repo.CreateUser(context.Background(),
		s.userID, fmt.Sprintf(models2.CommonFormat, "diff_users", s.userID))
	if err != nil {
		log.Println("failed to add user to list of active chat diff users", s.userID)
		s.notifyPeer(models2.Failed, map[string]any{
			"description": retryMessage,
		})
		s.peer.Close()
		return
	}

	/*
		this go-routine will exit when:
		(1) the user disconnects from chat manually
		(2) the app is closed
	*/
	go func() {
		log.Println("user joined", s.userID)
		ctx, cancel := context.WithCancel(context.Background())
		for {
			_, bMsg, err := s.peer.ReadMessage()
			if err != nil {
				_, ok := err.(*websocket.CloseError)
				if ok {
					log.Println("connection closed by user")
					s.disconnect()
				}
				cancel()
				return
			}

			cancel()
			ctx, cancel = context.WithCancel(context.Background())

			var chatList []models.Chat
			err = json.Unmarshal(bMsg, &chatList)
			if err != nil {
				continue
			}

			msgChan, err := s.diffRepo.SubscribeToChats(ctx, chatList)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			go func() {
				for {
					select {
					case msg := <-msgChan:
						fmt.Printf("got your message: %s from %s\n", msg.Payload, msg.SenderID.String())
						bMsg, err := json.Marshal(&msg)
						if err != nil {
							log.Println("failed to write message", err)
							continue
						}
						err = s.peer.WriteMessage(websocket.TextMessage, bMsg)
						if err != nil {
							log.Println("failed to write message", err)
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
	}()
}

func (s *DiffSession) notifyPeer(statusType models2.ConnectionStatusType,
	properties map[string]any) {
	status := models2.PeerConnectionStatus{
		Status:     statusType,
		Properties: properties,
	}
	notification := models2.Notification{
		Type: models2.InfoMessage,
		Body: &status,
	}
	bNotification, _ := json.Marshal(notification)
	err := s.peer.WriteMessage(websocket.TextMessage, bNotification)
	if err != nil {
		log.Println("failed to write message", err)
	}
}

// Invoked when the user disconnects (websocket connection is closed). It performs cleanup activities
func (s *DiffSession) disconnect() {
	//remove user from SET
	s.repo.RemoveUser(context.Background(),
		s.userID, fmt.Sprintf(models2.CommonFormat, "diff_users", s.userID))

	//close websocket
	s.peer.Close()
}

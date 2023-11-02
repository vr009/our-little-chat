package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/exp/slog"
	"log"
	"our-little-chatik/internal/models"
)

type PeerRepository struct {
	cl *redis.Client
}

func NewPeerRepository(cl *redis.Client) *PeerRepository {
	return &PeerRepository{
		cl: cl,
	}
}

func parseMessage(msgStr string) (*models.Message, error) {
	msg := models.Message{}
	err := json.Unmarshal([]byte(msgStr), &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *PeerRepository) SubscribeOnChatMessages(ctx context.Context,
	chatChannel string, readyChan chan struct{}) chan models.Message {
	/*
		this goroutine exits when the application shuts down. When the pusub connection is closed,
		the channel range loop terminates, hence terminating the goroutine
	*/
	messageChan := make(chan models.Message)
	go func() {
		log.Println("starting subscriber...", chatChannel)
		sub := r.cl.Subscribe(chatChannel)
		messages := sub.Channel()

		readyChan <- struct{}{}
		log.Println("LISTENING")
		for message := range messages {
			msg, err := parseMessage(message.Payload)
			log.Println("got one", msg.Payload)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			messageChan <- *msg
		}
		log.Println("SUBSCRIBER IS DOWN")
		select {
		case <-ctx.Done():
			err := sub.Unsubscribe(chatChannel)
			if err != nil {
				slog.Error(err.Error())
			}
			return
		default:
		}
	}()
	return messageChan
}

// SendMessageToChannel pusblishes on a redis pubsub channel
func (r *PeerRepository) SendMessageToChannel(ctx context.Context,
	msg models.Message, chatChannel string) {
	log.Println(msg.Payload, "sent to ", chatChannel)
	bMsg, err := json.Marshal(&msg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = r.cl.Publish(chatChannel, string(bMsg)).Err()
	if err != nil {
		log.Println("could not publish to channel", err)
	}
}

// CheckUserExists checks whether the user exists in the SET of active chat users
func (r *PeerRepository) CheckUserExists(ctx context.Context,
	user string, userSet string) (bool, error) {
	usernameTaken, err := r.cl.SIsMember(userSet, user).Result()
	if err != nil {
		return false, err
	}
	return usernameTaken, nil
}

// CreateUser creates a new user in the SET of active chat users
func (r *PeerRepository) CreateUser(ctx context.Context,
	user string, userSet string) error {
	err := r.cl.SAdd(userSet, user).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *PeerRepository) RemoveUser(ctx context.Context,
	user string, userSet string) {
	err := r.cl.SRem(userSet, user).Err()
	if err != nil {
		log.Println("failed to remove user:", user)
		return
	}
	log.Println("removed user from redis:", user)
}

func (r *PeerRepository) SaveMessage(message models.Message) error {
	bMsg, err := json.Marshal(&message)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", message.ChatID.String(), message.MsgID.String())
	err = r.cl.Set(key, string(bMsg), 0).Err()
	if err != nil {
		return err
	}
	key = fmt.Sprintf("%s:last-inserted", message.ChatID.String())
	err = r.cl.Set(key, string(bMsg), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

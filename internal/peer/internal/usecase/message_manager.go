package usecase

import (
	"container/list"
	"fmt"

	"our-little-chatik/internal/peer/internal"
	"our-little-chatik/internal/peer/internal/models"
)

type MessageManagerImpl struct {
	chatList *list.List
	repo     internal.PeerRepo
}

func NewMessageManager(repo internal.PeerRepo) *MessageManagerImpl {
	chatList := list.New()
	return &MessageManagerImpl{repo: repo, chatList: chatList}
}

func (m *MessageManagerImpl) EnqueueChat(chat *models.Chat) {
	ptr := chat
	m.chatList.PushFront(ptr)
}

func (m *MessageManagerImpl) DequeueChat(chat *models.Chat) {
	for e := m.chatList.Front(); e != nil; e = e.Next() {
		ch := e.Value.(*models.Chat)
		id := ch.ChatID
		if id == chat.ChatID {
			m.chatList.Remove(e)
		}
	}
}

// Work iterates all chats and checks all connected users.
// If someone has messages to receive it will put rhose messages to the channel.
// If someone has put a message to send, it will put it to a channel for sends.
func (m *MessageManagerImpl) Work() {
	for {
		for e := m.chatList.Front(); e != nil; e = e.Next() {
			chat := e.Value.(*models.Chat)
			for _, peer := range chat.Peers {
				if peer.Connected {
					msgs, _ := m.repo.FetchUpdates(chat, peer)
					if msgs != nil {
						peer.MsgsToRecv <- msgs
					}
					select {
					case msg := <-peer.MsgToSend:
						fmt.Println("Sending", msg)
						m.repo.SendPayload(msg)
					default:
					}
				}
			}
		}
	}
}

// EnqueueChatIfNotExists enqueues a passed Chat to an internal queue of chats.
// If the chat already exists it finds it and return.
func (m *MessageManagerImpl) EnqueueChatIfNotExists(c *models.Chat) (chat *models.Chat) {
	for e := m.chatList.Front(); e != nil; e = e.Next() {
		chat = e.Value.(*models.Chat)
		if chat.ChatID == c.ChatID {
			return chat
		}
	}
	chat = c
	m.EnqueueChat(chat)
	return chat
}

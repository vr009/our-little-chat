package usecase

import (
	"container/list"
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

func (m *MessageManagerImpl) Work() {
	for {
		for e := m.chatList.Front(); e != nil; e = e.Next() {
			chat := e.Value.(*models.Chat)
			msgs, _ := m.repo.FetchUpdates(chat)
			if msgs != nil {
				chat.PutMsgsToRecv(msgs)
			}
			select {
			case msg := <-chat.ReadyForSend:
				m.repo.SendPayload(msg, chat)
			default:
			}
		}
	}
}

// EnqueueChatIfNotExists enqueues a passed Chat to an internal queue of chats.
// If the chat already exists it finds it and return.
func (m *MessageManagerImpl) EnqueueChatIfNotExists(c *models.Chat) (chat *models.Chat) {
	for e := m.chatList.Front(); e != nil; e = e.Next() {
		chat = e.Value.(*models.Chat)
		if chat.ChatID == c.ChatID && chat.ReceiverID == c.ReceiverID {
			return chat
		}
	}
	chat = c
	m.EnqueueChat(chat)
	return chat
}

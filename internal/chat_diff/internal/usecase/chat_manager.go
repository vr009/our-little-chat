package usecase

import (
	"container/list"
	"fmt"
	"time"

	"our-little-chatik/internal/chat_diff/internal"
	"our-little-chatik/internal/models"
)

type ChatManager struct {
	repo  internal.ChatDiffRepo
	users *list.List
}

func NewChatManager(repo internal.ChatDiffRepo) *ChatManager {
	m := &ChatManager{repo: repo}
	m.users = list.New()
	return m
}

func (manager *ChatManager) AddChatUser(user *models.User) *models.User {
	for el := manager.users.Front(); el != manager.users.Back(); el = el.Next() {
		userToCompare := el.Value.(*models.User)
		if userToCompare.UserID == user.UserID {
			return userToCompare
		}
	}
	el := manager.users.PushBack(user)
	u := el.Value.(*models.User)
	fmt.Println("inserted", &u)
	return u
}

func (manager *ChatManager) Work() {
	for {
		for el := manager.users.Front(); el != nil; el = el.Next() {
			if el == nil {
				time.Sleep(time.Second)
				continue
			}
			user := el.Value.(*models.User)
			updates := manager.repo.FetchUpdates(*user)
			if len(updates) == 0 {
				continue
			}
			if updates != nil {
				fmt.Println("put to chan", &user)
				user.Updates <- updates
			}
		}
	}
}

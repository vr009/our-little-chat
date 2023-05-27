package usecase

import (
	"container/list"
	"fmt"
	"time"

	"our-little-chatik/internal/chat_diff/internal"
	models2 "our-little-chatik/internal/chat_diff/internal/models"
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

func (manager *ChatManager) AddChatUser(user *models2.ChatDiffUser) *models2.ChatDiffUser {
	for el := manager.users.Front(); el != manager.users.Back(); el = el.Next() {
		userToCompare := el.Value.(*models2.ChatDiffUser)
		if userToCompare.User.UserID == user.User.UserID {
			return user
		}
	}
	user.Updates = make(chan []models.ChatItem)
	el := manager.users.PushBack(user)
	u := el.Value.(*models2.ChatDiffUser)
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
			user := el.Value.(*models2.ChatDiffUser)
			updates := manager.repo.FetchUpdates(user.User)
			if len(updates) == 0 {
				continue
			}
			if updates != nil {
				go func() {
					fmt.Println("put to chan", &user)
					user.Updates <- updates
				}()
			}
		}
	}
}

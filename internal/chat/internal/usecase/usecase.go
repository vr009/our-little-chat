package usecase

import (
	"context"
	"fmt"
	"golang.org/x/exp/slices"
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type ChatUseCase struct {
	repo  internal.ChatRepo
	queue internal.QueueRepo
	users internal.UserDataInteractor
}

func NewChatUseCase(rep internal.ChatRepo, queue internal.QueueRepo,
	usersConnector internal.UserDataInteractor) *ChatUseCase {
	return &ChatUseCase{repo: rep, queue: queue, users: usersConnector}
}

func (ch *ChatUseCase) GetChatMessages(ctx context.Context, chat models.Chat,
	opts models.Opts) (models.Messages, models.StatusCode) {
	msgs, status := ch.queue.GetChatMessages(chat, opts)
	if status != models.OK {
		slog.Error("failed to fetch messages from queue %d", status)
	}
	if len(msgs) < int(opts.Limit) {
		opts.Limit = opts.Limit - int64(len(msgs))
		oldMsgs, status := ch.repo.GetChatMessages(ctx, chat, opts)
		if status != models.OK {
			slog.Error("failed to fetch messages from repo %d", status)
		}
		msgs = append(msgs, oldMsgs...)
	}
	return msgs, models.OK
}

func (ch *ChatUseCase) GetChatList(ctx context.Context, user models.User) ([]models.Chat, models.StatusCode) {
	chatList, statusCode := ch.repo.FetchChatList(ctx, user)
	if statusCode != models.OK {
		return nil, statusCode
	}
	lastChatsMessages, statusCode := ch.queue.GetChatsLastMessages(chatList)
	if statusCode == models.OK {
		for _, chat := range chatList {
			for _, msg := range lastChatsMessages {
				if chat.ChatID == msg.ChatID {
					chat.LastMessage = msg
				}
			}
		}
	}
	return chatList, models.OK
}

const defaultPhotoURL = "default.png"

func (ch *ChatUseCase) CreateChat(ctx context.Context, request models2.CreateChatRequest) (models.Chat, models.StatusCode) {
	chat := models.Chat{
		ChatID:       uuid.New(),
		CreatedAt:    time.Now().Unix(),
		Participants: request.Participants,
		PhotoURL:     *request.PhotoURL,
		Name:         *request.Name,
	}

	includeSelf := true
	for _, participant := range request.Participants {
		if participant == uuid.Nil {
			err := fmt.Errorf("added unexisting participant")
			slog.Error(err.Error())
			if err != nil {
				return models.Chat{}, models.BadRequest
			}
		}
		if participant == request.IssuerID {
			includeSelf = false
		}
	}
	if includeSelf {
		slog.Info("Adding a participant "+request.IssuerID.String()+" in list", "list",
			request.Participants)
		chat.Participants = append(request.Participants, request.IssuerID)
	}

	chatName := make(map[string]string)
	if len(chat.Participants) == 2 {
		for i := range chat.Participants {
			user, status := ch.users.GetUser(models.User{
				ID: chat.Participants[(i+1)%2],
			})
			if status != models.OK {
				chatName[chat.Participants[i].String()] = chat.ChatID.String()
			} else {
				chatName[chat.Participants[i].String()] = user.Nickname
			}
		}
	} else if len(chat.Participants) == 1 {
		user, status := ch.users.GetUser(models.User{
			ID: chat.Participants[0],
		})
		if status != models.OK {
			chatName[chat.Participants[0].String()] = chat.ChatID.String()
		} else {
			chatName[chat.Participants[0].String()] = user.Nickname
		}
	} else {
		for _, participant := range chat.Participants {
			chatName[participant.String()] = "Group chat " + chat.ChatID.String()
		}
	}
	if request.PhotoURL == nil {
		chat.PhotoURL = defaultPhotoURL
	}
	status := ch.repo.CreateChat(ctx, chat, chatName)
	if status != models.OK {
		return models.Chat{}, status
	}
	return chat, models.OK
}

func (ch *ChatUseCase) RemoveUserFromChat(ctx context.Context,
	chat models.Chat, users ...models.User) models.StatusCode {
	return ch.repo.RemoveUserFromChat(ctx, chat, users...)
}

func (ch *ChatUseCase) AddUsersToChat(ctx context.Context,
	chat models.Chat, users ...models.User) models.StatusCode {
	if len(users) == 0 {
		return models.BadRequest
	}

	chatFullInfo, status := ch.repo.GetChat(ctx, chat)
	if status != models.OK {
		return status
	}

	if len(chatFullInfo.Participants) == 1 {
		err := fmt.Errorf("it is not allowed to add users to the chat")
		slog.Error(err.Error())
		return models.Forbidden
	}

	usersToAdd := make([]models.User, 0)
	for _, user := range users {
		if !slices.Contains(chatFullInfo.Participants, user.ID) {
			usersToAdd = append(usersToAdd, user)
		}
	}

	chatNames := make(map[string]string)
	for _, user := range usersToAdd {
		if chatFullInfo.Name != "" {
			chatNames[user.ID.String()] = chatFullInfo.Name
		} else {
			chatNames[user.ID.String()] = "group chat " + chat.ChatID.String()[len(chat.ChatID.String())-5:]
		}
	}
	return ch.repo.AddUsersToChat(ctx, chat, chatNames, usersToAdd...)
}

func (ch *ChatUseCase) UpdateChatPhotoURL(ctx context.Context, chat models.Chat,
	photoURL string) models.StatusCode {
	return ch.repo.UpdateChatPhotoURL(ctx, chat, photoURL)
}

func (ch *ChatUseCase) GetChat(ctx context.Context, chat models.Chat) (models.Chat, models.StatusCode) {
	return ch.repo.GetChat(ctx, chat)
}

func (ch *ChatUseCase) DeleteChat(ctx context.Context, chat models.Chat) models.StatusCode {
	return ch.repo.DeleteChat(ctx, chat)
}

func (ch *ChatUseCase) DeleteMessage(ctx context.Context, message models.Message) models.StatusCode {
	return ch.repo.DeleteMessage(ctx, message)
}

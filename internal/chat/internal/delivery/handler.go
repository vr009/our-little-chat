package delivery

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"strconv"
)

type ChatEchoHandler struct {
	usecase        internal.ChatUseCase
	userInteractor internal.UserDataInteractor
}

func NewChatEchoHandler(usecase internal.ChatUseCase,
	userInteractor internal.UserDataInteractor) *ChatEchoHandler {
	return &ChatEchoHandler{
		usecase:        usecase,
		userInteractor: userInteractor,
	}
}

func (ch *ChatEchoHandler) GetChat(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	idStr := c.QueryParam("chat_id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: "passed empty parameter"})
	}

	chatID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: "bad id format"})
	}
	chat := models.Chat{ChatID: chatID}

	chat, err = ch.usecase.GetChat(chat)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, &chat)
}

func (ch *ChatEchoHandler) GetChatMessages(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	idStr := c.QueryParam("chat_id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Msg: "passed empty parameter for chat_id",
		})
	}
	offsetStr := c.QueryParam("offset")
	if offsetStr == "" {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Msg: "passed empty parameter offset",
		})
	}
	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Msg: "passed empty parameter limit",
		})
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{
			Msg: "passed empty parameter offset",
		})
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{
			Msg: "passed empty parameter limit",
		})
	}

	opts := models.Opts{Limit: limit, Page: offset}
	chatID, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Msg: "bad id format",
		})
	}

	chat := models.Chat{ChatID: chatID}
	msgs, err := ch.usecase.GetChatMessages(chat, opts)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, &models.Error{
			Msg: "internal issue",
		})
	}

	return c.JSON(http.StatusOK, &msgs)
}

func (ch *ChatEchoHandler) GetChatList(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	userID := c.Get("user_id").(uuid.UUID)
	user := models.User{ID: userID}

	chats, err := ch.usecase.GetChatList(user)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.Error{Msg: err.Error()})
	}

	log.Println("=============", chats)

	return c.JSON(http.StatusOK, &chats)
}

func (ch *ChatEchoHandler) PostNewChat(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	userID := c.Get("user_id").(uuid.UUID)
	user := models.User{ID: userID}

	chat := models.Chat{}
	err = c.Bind(&chat)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Msg: "bad body"})
	}

	log.Println("participants", chat.Participants)
	includeSelf := true
	for _, participant := range chat.Participants {
		if participant == uuid.Nil {
			errObj := models.Error{Msg: "Added unexisting participant"}
			if err != nil {
				return c.JSON(http.StatusBadRequest, errObj)
			}
		}
		if participant == user.ID {
			includeSelf = false
		}
	}
	if includeSelf {
		slog.Info("Adding a participant "+user.ID.String()+" in list", "list", chat.Participants)
		chat.Participants = append(chat.Participants, user.ID)
	}

	chatName := make(map[string]string)
	if len(chat.Participants) == 2 {
		for i := range chat.Participants {
			user, err := ch.userInteractor.GetUser(models.User{
				ID: chat.Participants[(i+1)%2],
			})
			if err != nil {
				//TODO chat id is not defined here
				chatName[chat.Participants[i].String()] = chat.ChatID.String()
			} else {
				chatName[chat.Participants[i].String()] = user.Nickname
			}
		}
	} else if len(chat.Participants) == 1 {
		user, err := ch.userInteractor.GetUser(models.User{
			ID: chat.Participants[0],
		})
		if err != nil {
			chatName[chat.Participants[0].String()] = chat.ChatID.String()
		} else {
			chatName[chat.Participants[0].String()] = user.Nickname
		}
	} else {
		for _, participant := range chat.Participants {
			chatName[participant.String()] = "Group chat " + chat.ChatID.String()
		}
	}
	log.Println(chat.Participants)
	log.Println("==== creating ==== ", chatName)
	createdChat, err := ch.usecase.CreateChat(chat, chatName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Msg: err.Error()})
	}
	createdChat.Name = chatName[user.ID.String()]

	return c.JSON(http.StatusCreated, &createdChat)
}

func (ch *ChatEchoHandler) AddUsersToChat(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	chat := models.Chat{}
	err = c.Bind(&chat)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Msg: "bad body"})
	}

	err = ch.usecase.UpdateChat(chat, models2.UpdateOptions{Action: models2.AddUsersToParticipants})
	if err != nil {
		errObj := models.Error{Msg: "Failed to add users"}
		return c.JSON(http.StatusBadRequest, errObj)
	}

	return c.JSON(http.StatusOK, models.Error{Msg: "OK"})
}

func (ch *ChatEchoHandler) RemoveUsersFromChat(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	chat := models.Chat{}
	err = c.Bind(&chat)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Msg: "bad body"})
	}

	err = ch.usecase.UpdateChat(chat, models2.UpdateOptions{Action: models2.RemoveUsersFromParticipants})
	if err != nil {
		errObj := models.Error{Msg: "Failed to add users"}
		return c.JSON(http.StatusBadRequest, errObj)
	}

	return c.JSON(http.StatusOK, models.Error{Msg: "OK"})
}

func (ch *ChatEchoHandler) ChangeChatPhoto(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	chat := models.Chat{}
	err = c.Bind(&chat)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Msg: "bad body"})
	}

	err = ch.usecase.UpdateChat(chat, models2.UpdateOptions{Action: models2.UpdatePhotoURL})
	if err != nil {
		errObj := models.Error{Msg: "Failed to update photo url"}
		return c.JSON(http.StatusBadRequest, errObj)
	}
	return c.JSON(http.StatusOK, models.Error{Msg: "OK"})
}

func (ch *ChatEchoHandler) DeleteChat(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	idStr := c.QueryParam("chat_id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: "passed empty parameter"})
	}

	chatID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: "bad id format"})
	}
	chat := models.Chat{ChatID: chatID}

	err = ch.usecase.DeleteChat(chat)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, &chat)
}

func (ch *ChatEchoHandler) DeleteMessage(c echo.Context) error {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	idStr := c.QueryParam("msg_id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: "passed empty parameter"})
	}

	msgID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: "bad id format"})
	}
	msg := models.Message{MsgID: msgID}

	err = ch.usecase.DeleteMessage(msg)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, &models.Error{Msg: err.Error()})
	}

	return c.JSON(http.StatusOK, &models.Error{Msg: "OK"})
}

package delivery

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	models2 "our-little-chatik/internal/chat/internal/models"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"our-little-chatik/internal/chat/internal"
	"our-little-chatik/internal/models"
)

type ChatEchoHandler struct {
	usecase internal.ChatUseCase
}

func NewChatEchoHandler(usecase internal.ChatUseCase) *ChatEchoHandler {
	return &ChatEchoHandler{
		usecase: usecase,
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
	user := models.User{UserID: userID}

	chats, err := ch.usecase.GetChatList(user)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, models.Error{Msg: err.Error()})
	}

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
	user := models.User{UserID: userID}

	chat := models.Chat{}
	err = c.Bind(&chat)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Msg: "bad body"})
	}

	includeSelf := true
	for _, participant := range chat.Participants {
		if participant == uuid.Nil {
			errObj := models.Error{Msg: "Added unexisting participant"}
			if err != nil {
				return c.JSON(http.StatusBadRequest, errObj)
			}
		}
		if participant == user.UserID {
			includeSelf = false
		}
	}
	if includeSelf {
		slog.Info("Adding a participant "+user.UserID.String()+" in list", "list", chat.Participants)
		chat.Participants = append(chat.Participants, user.UserID)
	}

	createdChat, err := ch.usecase.CreateChat(chat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Msg: err.Error()})
	}

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

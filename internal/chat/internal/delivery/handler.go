package delivery

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"net/http"
	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/internal/models"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/validator"
	"strconv"
	"time"
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

	v := validator.New()
	idStr := c.Param("id")
	v.Check(idStr != "", "id", "passed empty parameter")

	chatID, err := uuid.Parse(idStr)
	v.Check(err == nil, "id", "must be a correct uuid string value")
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	chat := models.Chat{ChatID: chatID}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var status models.StatusCode
	chat, status = ch.usecase.GetChat(ctx, chat)
	if status != models.OK {
		return pkg.ErrorResponse(c, http.StatusBadRequest, err.Error())
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

	v := validator.New()

	idStr := c.Param("id")
	v.Check(idStr != "", "id", "must be provided")

	offsetStr := c.QueryParam("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	v.Check(err == nil, "offset", "must be a correct integer value")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	v.Check(err == nil, "limit", "must be a correct integer value")

	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	opts := models.Opts{Limit: limit, Page: offset}
	chatID, err := uuid.Parse(idStr)
	if err != nil {
		return pkg.BadRequestResponse(c, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chat := models.Chat{ChatID: chatID}
	msgs, status := ch.usecase.GetChatMessages(ctx, chat, opts)
	if status != models.OK {
		return pkg.ErrorResponse(c, http.StatusInternalServerError, "internal issue")
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chats, status := ch.usecase.GetChatList(ctx, user)
	if status != models.OK {
		return pkg.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	input := models2.CreateChatRequest{}
	err = c.Bind(&input)
	if err != nil {
		return pkg.ErrorResponse(c, http.StatusBadRequest, "bad body")
	}

	v := validator.New()
	models2.ValidateLoginRequest(v, input)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	input.IssuerID = userID
	createdChat, status := ch.usecase.CreateChat(ctx, input)
	if status != models.OK {
		return pkg.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	input := models2.AddUsersToChatRequest{}
	err = c.Bind(&input)
	if err != nil {
		return pkg.ErrorResponse(c, http.StatusBadRequest, "bad body")
	}

	v := validator.New()
	models2.ValidateAddUsersToChatRequest(v, input)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	users := make([]models.User, len(input.Participants))
	for i := range users {
		users[i] = models.User{
			ID: input.Participants[i],
		}
	}

	chat := models.Chat{ChatID: *input.ChatID}

	status := ch.usecase.AddUsersToChat(ctx, chat, users...)
	if status != models.OK {
		return pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to add users")
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

	input := models2.RemoveUsersFromChatRequest{}
	err = c.Bind(&input)
	if err != nil {
		return pkg.ErrorResponse(c, http.StatusBadRequest, "bad body")
	}

	v := validator.New()
	models2.ValidateRemoveUsersFromChatRequest(v, input)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	users := make([]models.User, len(input.Participants))
	for i := range users {
		users[i] = models.User{
			ID: input.Participants[i],
		}
	}

	chat := models.Chat{ChatID: *input.ChatID}

	status := ch.usecase.RemoveUserFromChat(ctx, chat, users...)
	if status != models.OK {
		return pkg.ErrorResponse(c, http.StatusBadRequest, "Failed to add users")
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

	input := models2.UpdateChatPhotoURLRequest{}
	err = c.Bind(&input)
	if err != nil {
		return pkg.ErrorResponse(c, http.StatusBadRequest, "bad body")
	}

	v := validator.New()
	models2.ValidateUpdateChatPhotoURLRequest(v, input)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chat := models.Chat{ChatID: *input.ChatID}

	status := ch.usecase.UpdateChatPhotoURL(ctx, chat, *input.PhotoURL)
	if status != models.OK {
		return pkg.ErrorResponse(c, http.StatusInternalServerError, "Failed to update photo url")
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

	v := validator.New()
	idStr := c.QueryParam("chat_id")
	v.Check(idStr != "", "chat_id", "must be provided")

	chatID, err := uuid.Parse(idStr)
	v.Check(err == nil, "chat_id", "must be a correct uuid value")

	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	chat := models.Chat{ChatID: chatID}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	status := ch.usecase.DeleteChat(ctx, chat)
	if status != models.Deleted {
		return pkg.ErrorResponse(c, http.StatusBadRequest, err.Error())
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

	v := validator.New()
	idStr := c.QueryParam("msg_id")
	v.Check(idStr == "", "msg_id", "must be provided")

	msgID, err := uuid.Parse(idStr)
	v.Check(err == nil, "msg_id", "must be a correct uuid value")

	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	msg := models.Message{MsgID: msgID}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	status := ch.usecase.DeleteMessage(ctx, msg)
	if status != models.Deleted {
		return pkg.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &models.Error{Msg: "OK"})
}

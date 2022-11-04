package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"

	"github.com/google/uuid"
)

type ChatHandler struct {
	usecase internal.ChatUseCase
}

func NewChatHandler(usecase internal.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
	}
}

// GetChatMessages godoc
// @Summary Fetch the chat
// @Description get chat by ID
// @Produce  json
// @Success 200 {object} []internal.models.Message
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat/conv [get]
func (c *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	idStr := r.Header.Get("CHAT_ID")
	offset, err := strconv.ParseInt(r.Header.Get("OFFSET"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limit, err := strconv.ParseInt(r.Header.Get("LIMIT"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	opts := models.Opts{Limit: limit, Page: offset}
	chatID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat := models2.Chat{ChatID: chatID}
	msgs, err := c.usecase.FetchChatMessages(chat, opts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(msgs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetChatList godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} []models.ChatItem
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /accounts/{id} [get]
func (clh *ChatHandler) GetChatList(w http.ResponseWriter, r *http.Request) {
	var err error
	idStr := r.Header.Get("UserID")
	user := models.User{}
	user.UserID, err = uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chats, err := clh.usecase.GetChatList(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(chats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func (clh *ChatHandler) PostNewChat(w http.ResponseWriter, r *http.Request) {
	var err error
	chat := models2.Chat{}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdChat, err := clh.usecase.CreateNewChat(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(createdChat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func (clh *ChatHandler) PostChat(w http.ResponseWriter, r *http.Request) {
	var err error
	chat := models2.Chat{}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = clh.usecase.ActivateChat(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

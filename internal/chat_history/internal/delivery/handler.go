package delivery

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"our-little-chatik/internal/chat_history/internal"
	"our-little-chatik/internal/models"
	"strconv"
)

type ChatHandler struct {
	Usecase internal.ChatUseCase
}

func NewChatHandler(usecase internal.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		Usecase: usecase,
	}
}

func (c *ChatHandler) PostMessages(w http.ResponseWriter, r *http.Request) {
	msgs := []models.Message{}

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&msgs)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := c.Usecase.SaveMessages(msgs); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.WriteHeader(http.StatusOK)
}

// GetChat godoc
// @Summary Fetch the chat
// @Description get chat by ID
// @Produce  json
// @Success 200 {object} []internal.models.Message
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat/conv [get]
func (c *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
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
	chat := models.Chat{ChatID: chatID}
	msgs, err := c.Usecase.FetchChat(chat, opts)
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

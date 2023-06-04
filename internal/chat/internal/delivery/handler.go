package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"our-little-chatik/internal/chat/internal"
	models2 "our-little-chatik/internal/chat/models"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

type ChatHandler struct {
	usecase internal.ChatUseCase
}

func NewChatHandler(usecase internal.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		usecase: usecase,
	}
}

func (c *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	_, err = pkg.AuthHook(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errObj := models.Error{Msg: "Invalid token"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}
	idStr := r.URL.Query().Get("chat_id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		errObj := models.Error{Msg: "passed empty parameter"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}

	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	opts := models.Opts{Limit: limit, Page: offset}
	chatID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat := models2.Chat{ChatID: chatID}
	msgs, err := c.usecase.FetchChatMessages(chat, opts)
	if err != nil {
		fmt.Println(err)
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

func (clh *ChatHandler) GetChatList(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	user, err := pkg.AuthHook(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errObj := models.Error{Msg: "Invalid token"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}

	chats, err := clh.usecase.GetChatList(*user)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(&chats)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func (clh *ChatHandler) PostNewChat(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	usr, err := pkg.AuthHook(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errObj := models.Error{Msg: "Invalid token"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}
	chat := models2.Chat{}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	includeSelf := true
	for _, participant := range chat.Participants {
		if participant == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
			errObj := models.Error{Msg: "Added unexisting participant"}
			body, _ := json.Marshal(errObj)
			w.Write(body)
			return
		}
		if participant == usr.UserID {
			includeSelf = false
		}
	}
	if includeSelf {
		slog.Info("Adding a participant "+usr.UserID.String()+" in list", "list", chat.Participants)
		chat.Participants = append(chat.Participants, usr.UserID)
	}

	createdChat, err := clh.usecase.CreateNewChat(chat)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(createdChat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (clh *ChatHandler) PostChat(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	usr, err := pkg.AuthHook(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errObj := models.Error{Msg: "Invalid token"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}
	chat := models2.Chat{}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chat.Participants = append(chat.Participants, usr.UserID)
	err = clh.usecase.ActivateChat(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

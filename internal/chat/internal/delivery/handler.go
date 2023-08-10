package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	models2 "our-little-chatik/internal/chat/internal/models"
	"strconv"

	"our-little-chatik/internal/chat/internal"
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

func (c *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
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

	chatID, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat := models.Chat{ChatID: chatID}

	chat, err = c.usecase.GetChat(chat)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("FETCHED CHAT!!!!", chat)

	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(chat)
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
	chat := models.Chat{ChatID: chatID}
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

func (h *ChatHandler) GetChatList(w http.ResponseWriter, r *http.Request) {
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

	chats, err := h.usecase.GetChatList(*user)
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

func (h *ChatHandler) PostNewChat(w http.ResponseWriter, r *http.Request) {
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
	chat := models.Chat{}

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

	createdChat, err := h.usecase.CreateNewChat(chat)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("CREATED CHAT !!!!!!!!!!! ", createdChat)

	body, err := json.Marshal(createdChat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (h *ChatHandler) AddUsersToChat(w http.ResponseWriter, r *http.Request) {
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
	chat := models.Chat{}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecase.UpdateChat(chat, models2.UpdateOptions{Action: models2.AddUsersToParticipants})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errObj := models.Error{Msg: "Failed to add users"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) RemoveUsersFromChat(w http.ResponseWriter, r *http.Request) {
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
	chat := models.Chat{}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecase.UpdateChat(chat, models2.UpdateOptions{Action: models2.RemoveUsersFromParticipants})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errObj := models.Error{Msg: "Failed to add users"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) ChangeChatPhoto(w http.ResponseWriter, r *http.Request) {
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
	chat := models.Chat{}

	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecase.UpdateChat(chat, models2.UpdateOptions{Action: models2.UpdatePhotoURL})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errObj := models.Error{Msg: "Failed to update photo url"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}
	w.WriteHeader(http.StatusOK)
}

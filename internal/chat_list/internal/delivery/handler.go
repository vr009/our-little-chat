package delivery

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"our-little-chatik/internal/chat_list/internal"
	"our-little-chatik/internal/models"
)

type ChatListHandler struct {
	useacse internal.Usecase
}

func NewChatListHandler(useacse internal.Usecase) *ChatListHandler {
	return &ChatListHandler{useacse: useacse}
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
func (clh *ChatListHandler) GetChatList(w http.ResponseWriter, r *http.Request) {
	var err error
	idStr := r.Header.Get("UserID")
	user := models.User{}
	user.UserID, err = uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chats, err := clh.useacse.GetChatList(user)
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

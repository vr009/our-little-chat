package delivery

import (
	"encoding/json"
	"net/http"

	"our-little-chatik/internal/gateway/internal"
	"our-little-chatik/internal/models"

	"github.com/google/uuid"
)

type GatewayHandler struct {
	uc internal.GatewayUsecase
}

func NewGatewayHandler(uc internal.GatewayUsecase) *GatewayHandler {
	return &GatewayHandler{
		uc: uc,
	}
}

func (h *GatewayHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	session, err := h.uc.SignIn(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Set-Cookie", session.Token)
}

func (h *GatewayHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	session, err := h.uc.SignUp(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Set-Cookie", session.Token)
}

func (h *GatewayHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	var err error
	user := models.User{}
	idStr := r.Header.Get("UserID")
	user.UserID, err = uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	session, err := h.uc.SignUp(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	body, err := json.Marshal(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *GatewayHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	session := models.Session{}
	session.Token = r.Header.Get("Token")

	user, err := h.uc.GetUserFromSession(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

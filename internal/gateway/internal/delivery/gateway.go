package delivery

import (
	"encoding/json"
	"net/http"

	"our-little-chatik/internal/gateway/internal"
	"our-little-chatik/internal/models"
)

type GatewayHandler struct {
	uc internal.GatewayUsecase
}

func NewGatewayHandler(uc internal.GatewayUsecase) *GatewayHandler {
	return &GatewayHandler{
		uc: uc,
	}
}

type Error struct {
	Msg string `json:"message"`
}

func (h *GatewayHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errObj := &Error{Msg: "Bad body"}
		body, _ := json.Marshal(&errObj)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	session, err := h.uc.SignIn(&user)
	if err != nil {
		errObj := &Error{Msg: "Failed to sign in"}
		body, _ := json.Marshal(&errObj)
		w.WriteHeader(http.StatusForbidden)
		w.Write(body)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "Token", Value: session.Token})
	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := h.uc.SignUp(&user)
	if err != nil {
		errObj := &Error{Msg: "Failed to sign up"}
		body, _ := json.Marshal(&errObj)
		w.WriteHeader(http.StatusForbidden)
		w.Write(body)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "Token", Value: session.Token})
	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	session := models.Session{}
	cookies := r.Cookies()
	for i := range cookies {
		if cookies[i].Name == "Token" {
			session.Token = cookies[0].Value
		}
	}
	if session.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.uc.LogOut(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
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

// SignIn godoc
// @Summary Sign in a user
// @Accept  json
// @Success 200
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 500 {object} Error
// @Router /api/gateway/signin [post]
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

// SignUp godoc
// @Summary Sign up a user
// @Accept  json
// @Success 200
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 500 {object} Error
// @Router /api/gateway/signup [post]
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

// LogOut godoc
// @Summary Log out a user
// @Accept  json
// @Success 200
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /api/gateway/logout [delete]
func (h *GatewayHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	session := models.Session{}
	cookie, err := r.Cookie("Token")
	session.Token = cookie.Value
	if session.Token == "" {
		errObj := &Error{Msg: "No cookie provided"}
		body, _ := json.Marshal(&errObj)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	err = h.uc.LogOut(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) Find(w http.ResponseWriter, r *http.Request) {
	session := models.Session{}
	cookie, err := r.Cookie("Token")
	session.Token = cookie.Value
	if session.Token == "" {
		errObj := &Error{Msg: "No cookie provided"}
		body, _ := json.Marshal(&errObj)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	name := r.URL.Query().Get("name")

	users, err := h.uc.FindUser(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(&users)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

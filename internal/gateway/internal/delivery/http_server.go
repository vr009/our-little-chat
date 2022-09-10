package delivery

import (
	"gateway/internal"
	"net/http"
)

type Handler struct {
	uc internal.ServerUsecase
}

func NewHandler(uc internal.ServerUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

}

package delivery

import (
	"auth/internal"
	"log"
	"net/http"
)

type AuthHandler struct {
	useCase internal.AuthUseCase
}

func NewAuthHandler(useCase internal.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

func (ah *AuthHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	log.Print("GetSession")

	token := r.Header.Get("Token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
}

func (ah *AuthHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	log.Fatalf("DeleteSession")
}

func (ah *AuthHandler) PostSession(w http.ResponseWriter, r *http.Request) {
	log.Fatalf("PostSession")
}

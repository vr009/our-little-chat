package delivery

import (
	"net/http"
	"user_data/internal"
)

type UserdataHandler struct {
	useCase internal.UserdataUseCase
}

func NewUserdataHandler(useCase internal.UserdataUseCase) *UserdataHandler {
	return &UserdataHandler{
		useCase: useCase,
	}
}

func (udh *UserdataHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	return
}

func (udh *UserdataHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	return
}

func (udh *UserdataHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	return
}

func (udh *UserdataHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	return
}

func (udh *UserdataHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	return
}

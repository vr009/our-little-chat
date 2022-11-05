package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"

	"github.com/google/uuid"
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
	users, status := udh.useCase.GetAllUsers()
	if status == models.OK {
		body, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)

		if err != nil {
			log.Print("error")
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (udh *UserdataHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		log.Print("error of decoding")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, errCode := udh.useCase.CreateUser(person)
	fmt.Println("HERE", errCode)
	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	buf, err := json.Marshal(&s)
	if err != nil {
		log.Print("error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf)
	if err != nil {
		log.Print("error")
		return
	}
}

func (udh *UserdataHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Print("Get user")

	userID := r.Header.Get("user_id")

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuidFormString, err := uuid.Parse(userID)

	if err != nil {
		log.Print("error of UUID parsing")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person := models.UserData{
		UserID: uuidFormString,
	}

	s, errCode := udh.useCase.GetUser(person)

	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	a, err := json.Marshal(&s)

	if err != nil {
		log.Print("error of json.Marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(a)

	if err != nil {
		log.Print("error")
		return
	}
}

func (udh *UserdataHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		log.Print("error of decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errCode := udh.useCase.DeleteUser(person)

	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (udh *UserdataHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		log.Print("error of decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, errCode := udh.useCase.UpdateUser(person)

	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	a, err := json.Marshal(&s)

	if err != nil {
		log.Print("error of json.Marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(a)

	if err != nil {
		log.Print("error")
		return
	}

}

func (udh *UserdataHandler) CheckUserData(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		log.Print("error of decoding")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errCode := udh.useCase.CheckUser(person)

	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleErrorCode(errCode models.StatusCode, w http.ResponseWriter) {
	if errCode == models.NotFound {
		log.Print("error of StatusNotFound")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if errCode == models.InternalError {
		log.Print("error of StatusInternalServerError")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if errCode == models.BadRequest {
		log.Print("error of StatusInternalServerError")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errCode != models.OK {
		log.Print("error of StatusInternalServerError")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

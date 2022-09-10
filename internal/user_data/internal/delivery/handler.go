package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"
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
	fmt.Println(users, status)
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

	checkErrorCode(errCode, w)

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

	checkErrorCode(errCode, w)

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
	log.Print("DeleteSession")

	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		log.Print("error of decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errCode := udh.useCase.DeleteUser(person)

	if errCode != models.OK {
		checkErrorCode(errCode, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (udh *UserdataHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Print("Update user")

	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		log.Print("error of decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, errCode := udh.useCase.UpdateUser(person)

	checkErrorCode(errCode, w)

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

func checkErrorCode(errCode models.StatusCode, w http.ResponseWriter) {
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

	if errCode != models.OK {
		log.Print("error of StatusInternalServerError")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

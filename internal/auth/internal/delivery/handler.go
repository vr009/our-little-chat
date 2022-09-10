package delivery

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"our-little-chatik/internal/auth/internal"
	"our-little-chatik/internal/auth/internal/models"
)

type AuthHandler struct {
	useCase internal.AuthUseCase
}

func NewAuthHandler(useCase internal.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

func (ah *AuthHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	log.Print("GetToken")

	userID := r.Header.Get("UserID")

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

	session := models.Session{
		UserID: uuidFormString,
	}

	s, errCode := ah.useCase.GetToken(session)

	CheckErrorCode(errCode, w)

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

func (ah *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Print("GetUser")

	token := r.Header.Get("Token")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session := models.Session{
		Token: token,
	}

	s, errCode := ah.useCase.GetUser(session)

	CheckErrorCode(errCode, w)

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

func (ah *AuthHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteSession")

	session := models.Session{}

	err := json.NewDecoder(r.Body).Decode(&session)

	if err != nil {
		log.Print("error of decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errCode := ah.useCase.DeleteSession(session)

	if errCode != models.OK {
		CheckErrorCode(errCode, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ah *AuthHandler) PostSession(w http.ResponseWriter, r *http.Request) {

	session := models.Session{}

	err := json.NewDecoder(r.Body).Decode(&session)

	if err != nil {
		log.Print("error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, errCode := ah.useCase.CreateSession(session)

	CheckErrorCode(errCode, w)

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

func CheckErrorCode(errCode models.StatusCode, w http.ResponseWriter) {
	switch errCode {
	case models.NotFound:
		{
			log.Print("error of StatusNotFound")
			w.WriteHeader(http.StatusNotFound)
		}
	case models.InternalError:
		{
			log.Print("error of StatusInternalServerError")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

}

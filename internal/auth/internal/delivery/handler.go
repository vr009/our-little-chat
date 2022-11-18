package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"our-little-chatik/internal/auth/internal"
	"our-little-chatik/internal/auth/internal/models"
	models2 "our-little-chatik/internal/models"

	"github.com/google/uuid"
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

	session := models2.Session{
		UserID: uuidFormString,
	}

	s, errCode := ah.useCase.GetToken(session)

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

func (ah *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Print("GetUser")

	token := r.Header.Get("Token")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session := models2.Session{
		Token: token,
	}

	s, errCode := ah.useCase.GetUser(session)

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

func (ah *AuthHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteSession")

	session := models2.Session{}

	err := json.NewDecoder(r.Body).Decode(&session)

	if err != nil {
		log.Print("error of decoding request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errCode := ah.useCase.DeleteSession(session)

	if errCode != models.OK {
		checkErrorCode(errCode, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ah *AuthHandler) PostSession(w http.ResponseWriter, r *http.Request) {
	session := models2.Session{}

	err := json.NewDecoder(r.Body).Decode(&session)

	if err != nil {
		log.Print("error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, errCode := ah.useCase.CreateSession(session)

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

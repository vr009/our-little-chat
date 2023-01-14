package delivery

import (
	"encoding/json"
	"net/http"

	"our-little-chatik/internal/auth/internal"
	"our-little-chatik/internal/auth/internal/models"
	models2 "our-little-chatik/internal/models"

	"github.com/golang/glog"
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
	userID := r.Header.Get("UserID")

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuidFormString, err := uuid.Parse(userID)

	if err != nil {
		glog.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session := models2.Session{
		UserID: uuidFormString,
	}

	s, errCode := ah.useCase.GetToken(session)

	if errCode != models.OK {
		checkErrorCode(errCode, w)
		return
	}

	a, err := json.Marshal(&s)

	if err != nil {
		glog.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(a)

	if err != nil {
		glog.Errorf(err.Error())
		return
	}
}

func (ah *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	glog.Warningf("hit %s", token)
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	glog.Warningf("ok %s", token)
	session := models2.Session{
		Token: token,
	}

	s, errCode := ah.useCase.GetUser(session)

	user := models2.User{UserID: s.UserID}

	if errCode != models.OK {
		glog.Warningf("not ok %s", token)
		checkErrorCode(errCode, w)
		return
	}

	body, err := json.Marshal(&user)
	glog.Warningf("%v", s.UserID)
	if err != nil {
		glog.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)

	if err != nil {
		glog.Errorf(err.Error())
		return
	}
	glog.Warningf("ok ok %v", s)
}

func (ah *AuthHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session := models2.Session{}

	session.Token = r.Header.Get("Token")
	if session.Token == "" {
		glog.Errorf("Received empty token")
		w.WriteHeader(http.StatusBadRequest)
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
		glog.Errorf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, errCode := ah.useCase.CreateSession(session)

	if errCode != models.OK {
		glog.Errorf("Failed to create a session")
		checkErrorCode(errCode, w)
		return
	}

	buf, err := json.Marshal(&s)
	if err != nil {
		glog.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf)
	if err != nil {
		glog.Errorf(err.Error())
		return
	}
}

func checkErrorCode(errCode models.StatusCode, w http.ResponseWriter) {
	if errCode == models.NotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if errCode == models.InternalError {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if errCode != models.OK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

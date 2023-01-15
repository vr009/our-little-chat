package delivery

import (
	"encoding/json"
	"net/http"

	"our-little-chatik/internal/common"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"

	"github.com/golang/glog"
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
			glog.Errorf(err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (udh *UserdataHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		glog.Errorf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	glog.Infof("Unmarshalled: %v", person)

	newPeson, errCode := udh.useCase.CreateUser(person)
	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	glog.Infof("Created: %v", newPeson)

	buf, err := json.Marshal(&newPeson)
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

var defaultAuthUrl = "http://auth:8087/api/v1/auth/user"

func (udh *UserdataHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := common.AuthHook(r, defaultAuthUrl)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errObj := models2.Error{Msg: "Invalid token"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}

	if err != nil {
		glog.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person := models.UserData{}
	person.UserID = user.UserID

	s, errCode := udh.useCase.GetUser(person)

	if errCode != models.OK {
		handleErrorCode(errCode, w)
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

func (udh *UserdataHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		glog.Errorf(err.Error())
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
		glog.Errorf(err.Error())
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

func (udh *UserdataHandler) CheckUserData(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		glog.Errorf(err.Error())
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

func (udh *UserdataHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	glog.Infof("Searching for %s", name)
	users, errCode := udh.useCase.FindUser(name)
	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}
	body, err := json.Marshal(users)
	if err != nil {
		glog.Errorf("Failed to marshal body for users: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func handleErrorCode(errCode models.StatusCode, w http.ResponseWriter) {
	if errCode == models.NotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if errCode == models.InternalError {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if errCode == models.BadRequest {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errCode != models.OK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

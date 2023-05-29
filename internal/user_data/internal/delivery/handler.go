package delivery

import (
	"encoding/json"
	"net/http"

	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"

	"golang.org/x/exp/slog"
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
			slog.Error(err.Error())
			return
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (udh *UserdataHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Info("Unmarshalled:", "person", slog.AnyValue(person))

	newPerson, errCode := udh.useCase.CreateUser(person)
	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	slog.Info("Created: ", slog.AnyValue(newPerson))

	buf, err := json.Marshal(&newPerson)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (udh *UserdataHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := pkg.AuthHook(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errObj := models2.Error{Msg: "Invalid token"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}

	if err != nil {
		slog.Error(err.Error())
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
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(a)

	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (udh *UserdataHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		slog.Error(err.Error())
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
		slog.Error(err.Error())
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
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(a)

	if err != nil {
		slog.Error(err.Error())
		return
	}

}

func (udh *UserdataHandler) CheckUserData(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, errCode := udh.useCase.CheckUser(person)

	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (udh *UserdataHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	_, err := pkg.AuthHook(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		errObj := models2.Error{Msg: "Invalid token"}
		body, _ := json.Marshal(errObj)
		w.Write(body)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	slog.Info("Searching for " + name)
	users, errCode := udh.useCase.FindUser(name)
	if errCode != models.OK {
		handleErrorCode(errCode, w)
		return
	}
	body, err := json.Marshal(users)
	if err != nil {
		slog.Error("Failed to marshal body for users: " + err.Error())
		handleErrorCode(models.InternalError, w)
		return
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

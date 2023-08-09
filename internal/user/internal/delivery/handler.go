package delivery

import (
	"encoding/json"
	"net/http"

	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/user/internal"
	"our-little-chatik/internal/user/internal/models"

	"github.com/google/uuid"
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
		handleErrorCode(errCode, w, models.Error{Msg: "Failed to create user"})
		return
	}

	slog.Info("Created: ", slog.AnyValue(newPerson))

	buf, err := json.Marshal(&newPerson)
	if err != nil {
		slog.Error(err.Error())
		handleErrorCode(errCode, w, models.Error{Msg: "Failed to marshal the response body"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(buf)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (udh *UserdataHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := pkg.AuthHook(r)
	if err != nil {
		handleErrorCode(models.Forbidden, w, models.Error{Msg: "Invalid token"})
		return
	}

	person := models.UserData{}
	person.UserID = user.UserID

	s, errCode := udh.useCase.GetUser(person)
	if errCode != models.OK {
		handleErrorCode(errCode, w, models.Error{Msg: "Failed to find info about the user"})
		return
	}

	a, err := json.Marshal(&s)
	if err != nil {
		slog.Error(err.Error())
		handleErrorCode(models.InternalError, w, models.Error{Msg: "Failed to marshal response body"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(a)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (udh *UserdataHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := pkg.AuthHook(r)
	if err != nil {
		handleErrorCode(models.Forbidden, w, models.Error{Msg: "Invalid token"})
		slog.Error(err.Error())
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleErrorCode(models.BadRequest, w, models.Error{Msg: "Bad id format"})
		return
	}

	slog.Info("requested info", "from user", user.UserID.String(), "about", idStr)

	person := models.UserData{}
	person.UserID = id

	person, errCode := udh.useCase.GetUser(person)
	if errCode != models.OK {
		handleErrorCode(errCode, w, models.Error{Msg: "Failed to get info about the user"})
		return
	}

	a, err := json.Marshal(&person)
	if err != nil {
		slog.Error(err.Error())
		handleErrorCode(models.InternalError, w, models.Error{Msg: "Failed to marshal response body"})
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
		handleErrorCode(errCode, w, models.Error{Msg: "Failed to delete user"})
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
		handleErrorCode(errCode, w, models.Error{Msg: "User update failed"})
		return
	}

	a, err := json.Marshal(&s)
	if err != nil {
		slog.Error(err.Error())
		handleErrorCode(models.InternalError, w, models.Error{Msg: "Failed to marshal response body"})
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
		handleErrorCode(errCode, w, models.Error{Msg: "User data is inaccessible"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (udh *UserdataHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	_, err := pkg.AuthHook(r)
	if err != nil {
		handleErrorCode(models.Forbidden, w, models.Error{Msg: "Invalid token"})
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	slog.Info("Searching for " + name)

	users, errCode := udh.useCase.FindUser(name)
	if errCode != models.OK {
		handleErrorCode(errCode, w, models.Error{Msg: "Inaccessible user data"})
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		slog.Error("Failed to marshal body for users: " + err.Error())
		handleErrorCode(models.InternalError, w, models.Error{Msg: "Failed to marshal response body"})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func handleErrorCode(errCode models.StatusCode, w http.ResponseWriter, errObj models.Error) {
	switch errCode {
	case models.NotFound:
		w.WriteHeader(http.StatusNotFound)
	case models.InternalError:
		w.WriteHeader(http.StatusInternalServerError)
	case models.BadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case models.Forbidden:
		w.WriteHeader(http.StatusForbidden)
	default:
		if errCode != models.OK {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	body, _ := json.Marshal(errObj)
	w.Write(body)
}

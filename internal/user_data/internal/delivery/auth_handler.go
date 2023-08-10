package delivery

import (
	"encoding/json"
	"net/http"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/user_data/internal"
	"our-little-chatik/internal/user_data/internal/models"

	"golang.org/x/exp/slog"
)

type AuthHandler struct {
	useCase internal.UserdataUseCase
}

func NewAuthHandler(useCase internal.UserdataUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		slog.Error(err.Error())
		handleErrorCode(models.BadRequest, w, models.Error{Msg: "Bad format"})
		return
	}

	slog.Info("Unmarshalled:", "person", slog.AnyValue(person))
	newPerson, errCode := h.useCase.CreateUser(person)
	if errCode != models.OK {
		handleErrorCode(errCode, w, models.Error{Msg: "Creating user failed"})
		return
	}

	token, err := pkg.GenerateJWTToken(newPerson.User, false)
	if err != nil {
		slog.Error(err.Error())
		handleErrorCode(models.InternalError, w, models.Error{Msg: "Fail while token generating"})
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "Token", Value: token, Path: "/"})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	person := models.UserData{}

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		slog.Error(err.Error())
		handleErrorCode(models.BadRequest, w, models.Error{Msg: "Bad format"})
		return
	}

	slog.Info("Unmarshalled:", "person", slog.AnyValue(person))

	usr, code := h.useCase.CheckUser(person)
	if code != models.OK {
		handleErrorCode(code, w, models.Error{Msg: "Inaccessible user data"})
		return
	}

	person.User.UserID = usr.UserID

	token, err := pkg.GenerateJWTToken(person.User, false)
	if err != nil {
		handleErrorCode(models.InternalError, w, models.Error{Msg: "Fail while token generating"})
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "Token", Value: token, Path: "/"})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user, err := pkg.AuthHook(r)
	if err != nil {
		errObj := &models.Error{Msg: err.Error()}
		body, _ := json.Marshal(&errObj)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	token, err := pkg.GenerateJWTToken(*user, true)
	if err != nil {
		handleErrorCode(models.InternalError, w, models.Error{Msg: "Fail while token generating"})
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "Token", Value: token, Path: "/"})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

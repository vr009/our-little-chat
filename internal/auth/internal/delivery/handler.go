package delivery

import (
	"auth/internal"
	models2 "auth/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type AuthHandler struct {
	UseCase internal.UseCase
}

func NewAuthHandler(UCase internal.UseCase) *AuthHandler {
	return &AuthHandler{
		UseCase: UCase,
	}
}

// SignUp godoc
// @Summary      SignUp user
// @Description  sign up for user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param User body models.UserCreate true "Login user"
// @Success      201  {object}  models.User
// @Failure      400  {object}  models.Error
// @Failure      404  {object}  models.Error
// @Failure      500  {object}  models.Error
// @Router       /auth/signup [post]
func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	r.WithContext(context.Background())
	userCreate := models2.UserCreate{}
	user := models2.User{}
	err := json.NewDecoder(r.Body).Decode(&userCreate)
	log.Println(err)
	if err != nil {
		response(w, models2.INTERNAL, nil)
		return
	}
	user.Password = userCreate.Password
	user.Username = userCreate.Username
	user.Lastname = userCreate.Lastname
	user.Firstname = userCreate.Firstname

	authedUsr, errCode := a.UseCase.SignUp(&user)
	if errCode != models2.OK {
		response(w, errCode, nil)
		return
	}
	body, err := json.Marshal(authedUsr)
	log.Println(err)
	if err != nil {
		response(w, models2.INTERNAL, nil)
		return
	}
	response(w, errCode, body)
}

// SignIn godoc
// @Summary      SignIn user
// @Description  for log in of user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param User body models.UserLogin true "Create user"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Error
// @Failure      404  {object}  models.Error
// @Failure      500  {object}  models.Error
// @Router       /auth/signin [post]
func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	userLogin := models2.UserLogin{}

	user := models2.User{}
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	log.Println(err)
	if err != nil {
		errBody, _ := json.Marshal(models2.Error{Message: "Internal error"})
		response(w, models2.INTERNAL, errBody)
		return
	}

	user.Password = userLogin.Password
	user.Username = userLogin.Username

	log.Println(user)
	authedUsr, errCode := a.UseCase.SignIn(&user)
	if errCode != models2.OK {
		errBody, _ := json.Marshal(models2.Error{Message: "Some error"})
		response(w, errCode, errBody)
		return
	}
	body, err := json.Marshal(authedUsr)
	log.Println(err)
	if err != nil {
		errBody, _ := json.Marshal(models2.Error{Message: "Internal error"})
		response(w, models2.INTERNAL, errBody)
		return
	}
	response(w, errCode, body)
}

func response(w http.ResponseWriter, errCode models2.ErrorCode, body []byte) {
	switch errCode {
	case models2.OK:
		w.WriteHeader(http.StatusOK)
	case models2.EXISTS:
		w.WriteHeader(http.StatusConflict)
	case models2.NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
	case models2.CREATED:
		w.WriteHeader(http.StatusCreated)
	case models2.INTERNAL:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(body)
}

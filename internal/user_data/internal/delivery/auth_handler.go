package delivery

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/validator"
	"our-little-chatik/internal/user_data/internal"

	"golang.org/x/exp/slog"
)

type AuthEchoHandler struct {
	useCase internal.UserUsecase
}

func NewAuthEchoHandler(useCase internal.UserUsecase) *AuthEchoHandler {
	return &AuthEchoHandler{
		useCase: useCase,
	}
}

func (h *AuthEchoHandler) SignUp(c echo.Context) error {
	var input struct {
		Nickname string `json:"nickname,omitempty"`
		Name     string `json:"name,omitempty"`
		Surname  string `json:"surname,omitempty"`
		Password string `json:"password,omitempty"`
	}

	err := c.Bind(&input)
	if err != nil {
		slog.Error(err.Error())
		return pkg.BadRequestResponse(c, err)
	}

	person := &models.User{
		Name:      input.Name,
		Nickname:  input.Nickname,
		Surname:   input.Surname,
		Activated: true,
	}

	// TODO move to usecase
	err = person.Password.Set(input.Password)
	if err != nil {
		return pkg.ServerErrorResponse(c, err)
	}

	v := validator.New()
	validator.ValidateUserBySignUp(v, person)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	slog.Debug("Unmarshalled:", "person", slog.AnyValue(person))
	newPerson, errCode := h.useCase.CreateUser(*person)
	if errCode != models.OK {
		switch errCode {
		case models.Conflict:
			v.AddError("credentials", "this user already exists")
			return pkg.FailedValidationResponse(c, v.Errors)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf("%s", err.Error()))
		}
	}

	token, err := pkg.GenerateJWTTokenV2(newPerson, false)
	if err != nil {
		slog.Error(err.Error())
		return pkg.ServerErrorResponse(c, err)
	}
	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthEchoHandler) Login(c echo.Context) error {
	var input struct {
		Nickname string `json:"nickname,omitempty"`
		Password string `json:"password,omitempty"`
	}

	err := c.Bind(&input)
	if err != nil {
		slog.Error(err.Error())
		return pkg.BadRequestResponse(c, err)
	}

	person := &models.User{
		Nickname: input.Nickname,
	}
	v := validator.New()
	validator.ValidateUserBySignIn(v, person)
	validator.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	slog.Debug("Unmarshalled:", "person", slog.AnyValue(person))

	foundPerson, code := h.useCase.CheckUser(*person)
	if code != models.OK {
		switch code {
		case models.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf("internal issue"))
		}
	}

	// TODO move to usecase
	match, err := foundPerson.Password.Matches(input.Password)
	if err != nil {
		return pkg.ServerErrorResponse(c, err)
	}
	if !match {
		return pkg.InvalidCredentialsResponse(c)
	}

	token, err := pkg.GenerateJWTTokenV2(foundPerson, false)
	if err != nil {
		slog.Error(err.Error())
		return pkg.ServerErrorResponse(c, err)
	}
	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthEchoHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	user := models.User{ID: userID}

	log.Println("LOGOUT", user.ID)

	token, err := pkg.GenerateJWTTokenV2(user, true)
	if err != nil {
		slog.Error(err.Error())
		return pkg.ServerErrorResponse(c, err)
	}

	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

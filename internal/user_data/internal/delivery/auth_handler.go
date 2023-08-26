package delivery

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/user_data/internal"

	"golang.org/x/exp/slog"
)

type AuthEchoHandler struct {
	useCase internal.UserdataUsecase
}

func NewAuthEchoHandler(useCase internal.UserdataUsecase) *AuthEchoHandler {
	return &AuthEchoHandler{
		useCase: useCase,
	}
}

func (h *AuthEchoHandler) SignUp(c echo.Context) error {
	log.Println("HERE!!!!!", c.Request().Header.Get("Content-Type"))
	person := models.UserData{}
	err := c.Bind(&person)
	if err != nil {
		slog.Error(err.Error())
		return pkg.HandleErrorCode(models.BadRequest, models.Error{Msg: "bad body"}, c)
	}

	slog.Info("Unmarshalled:", "person", slog.AnyValue(person))
	newPerson, errCode := h.useCase.CreateUser(person)
	if errCode != models.OK {
		return pkg.HandleErrorCode(errCode, models.Error{Msg: "Creating user failed"}, c)
	}

	token, err := pkg.GenerateJWTTokenV2(newPerson.User, false)
	if err != nil {
		slog.Error(err.Error())
		return pkg.HandleErrorCode(models.InternalError, models.Error{Msg: "Fail while token generating"}, c)
	}
	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthEchoHandler) Login(c echo.Context) error {
	person := models.UserData{}
	err := c.Bind(&person)
	if err != nil {
		slog.Error(err.Error())
		return pkg.HandleErrorCode(models.BadRequest, models.Error{Msg: "bad body"}, c)
	}

	slog.Info("Unmarshalled:", "person", slog.AnyValue(person))

	usr, code := h.useCase.CheckUser(person)
	if code != models.OK {
		return pkg.HandleErrorCode(code, models.Error{Msg: "Inaccessible user data"}, c)
	}

	person.User.UserID = usr.UserID

	token, err := pkg.GenerateJWTTokenV2(person.User, false)
	if err != nil {
		return pkg.HandleErrorCode(models.InternalError, models.Error{Msg: "Fail while token generating"}, c)
	}
	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthEchoHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	user := models.User{UserID: userID}

	log.Println("LOGOUT", user.UserID)

	token, err := pkg.GenerateJWTTokenV2(user, true)
	if err != nil {
		return pkg.HandleErrorCode(models.InternalError, models.Error{Msg: "Fail while token generating"}, c)
	}

	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

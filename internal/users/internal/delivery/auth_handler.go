package delivery

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/validator"
	"our-little-chatik/internal/users/internal"
	models2 "our-little-chatik/internal/users/internal/models"

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

// SignUp godoc
// @Summary Sign up a user.
// @Description sign up a user.
// @Accept json
// @Produce json
// @Tags auth
// @Param request body models.SignUpPersonRequest true "sign up user request"
// @Success 303
// @Failure 401 {object} models.HttpResponse
// @Failure 422 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /users/signup [post]
func (h *AuthEchoHandler) SignUp(c echo.Context) error {
	input := models2.SignUpPersonRequest{}
	err := c.Bind(&input)
	if err != nil {
		slog.Error(err.Error())
		return pkg.BadRequestResponse(c, err)
	}

	v := validator.New()
	models2.ValidateSignUpRequest(v, input)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	newPerson, errCode := h.useCase.SignUp(input)
	if errCode != models.OK {
		switch errCode {
		case models.Conflict:
			v.AddError("credentials", "this user already exists")
			return pkg.UnauthorizedResponse(c, fmt.Errorf("failed to registrate with passed credentials"))
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

// Login godoc
// @Summary log in a user.
// @Description log in a user.
// @Accept json
// @Produce json
// @Tags auth
// @Param request body models.LoginRequest true "log in request"
// @Success 303
// @Failure 401 {object} models.HttpResponse
// @Failure 422 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /users/login [post]
func (h *AuthEchoHandler) Login(c echo.Context) error {
	input := models2.LoginRequest{}
	err := c.Bind(&input)
	if err != nil {
		slog.Error(err.Error())
		return pkg.BadRequestResponse(c, err)
	}

	v := validator.New()
	models2.ValidateLoginRequest(v, input)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	foundUser, code := h.useCase.Login(input)
	if code != models.OK {
		switch code {
		case models.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.UnauthorizedResponse(c, fmt.Errorf("internal issue"))
		}
	}

	token, err := pkg.GenerateJWTTokenV2(foundUser, false)
	if err != nil {
		slog.Error(err.Error())
		return pkg.ServerErrorResponse(c, err)
	}
	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

// Logout godoc
// @Summary Log out a user.
// @Description log out a user.
// @Accept json
// @Produce json
// @Tags auth
// @Success 303
// @Failure 401 {object} models.HttpResponse
// @Failure 422 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /users/logout [post]
func (h *AuthEchoHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	user := models.User{ID: userID}

	token, err := pkg.GenerateJWTTokenV2(user, true)
	if err != nil {
		slog.Error(err.Error())
		return pkg.ServerErrorResponse(c, err)
	}

	c.SetCookie(&http.Cookie{Name: "Token", Value: token, Path: "/"})
	return c.Redirect(http.StatusSeeOther, "/")
}

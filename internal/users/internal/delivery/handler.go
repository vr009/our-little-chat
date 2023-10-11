package delivery

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"net/http"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/validator"
	"our-little-chatik/internal/users/internal"
	"our-little-chatik/internal/users/internal/models"
)

type UserEchoHandler struct {
	useCase internal.UserUsecase
}

func NewUserEchoHandler(useCase internal.UserUsecase) *UserEchoHandler {
	return &UserEchoHandler{
		useCase: useCase,
	}
}

func (udh *UserEchoHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)

	me, errCode := udh.useCase.GetUser(models.GetUserRequest{UserID: userID})
	if errCode != models2.OK {
		switch errCode {
		case models2.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf(""))
		}
	}

	return c.JSON(http.StatusOK, &me)
}

func (udh *UserEchoHandler) GetUserForID(c echo.Context) error {
	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return pkg.BadRequestResponse(c, err)
	}

	foundPerson, errCode := udh.useCase.GetUser(models.GetUserRequest{UserID: userID})
	if errCode != models2.OK {
		switch errCode {
		case models2.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf("faile to get a user due to internal issue"))
		}
	}
	return c.JSON(http.StatusOK, &foundPerson)
}

func (udh *UserEchoHandler) DeactivateUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return pkg.BadRequestResponse(c, fmt.Errorf("issuer info not provided"))
	}

	person := models2.User{
		ID: userID,
	}

	errCode := udh.useCase.DeactivateUser(person)
	if errCode != models2.OK {
		switch errCode {
		case models2.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf("internal issue"))
		}
	}
	return c.JSON(http.StatusOK, models2.Error{Msg: "OK"})
}

func (udh *UserEchoHandler) UpdateUser(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	user := models2.User{ID: userID}

	var input models.UpdateUserRequest
	err := c.Bind(&input)
	if err != nil {
		slog.Error(err.Error())
		return pkg.BadRequestResponse(c, err)
	}

	v := validator.New()
	models.ValidateUpdateUserRequest(v, input)

	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	s, errCode := udh.useCase.UpdateUser(user, input)
	if errCode != models2.OK {
		switch errCode {
		case models2.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf("internal issue"))
		}
	}

	return c.JSON(http.StatusOK, &s)
}

func (udh *UserEchoHandler) SearchUsers(c echo.Context) error {
	name := c.QueryParam("nickname")
	v := validator.New()
	v.Check(name != "", "nickname", "must be provided")
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}
	slog.Info("Searching for " + name)

	users, errCode := udh.useCase.FindUsers(name)
	if errCode != models2.OK {
		switch errCode {
		case models2.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf("internal issue"))
		}
	}

	return c.JSON(http.StatusOK, users)
}

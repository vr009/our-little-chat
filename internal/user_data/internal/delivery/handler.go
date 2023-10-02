package delivery

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/pkg/validator"
	"our-little-chatik/internal/user_data/internal"
)

type UserEchoHandler struct {
	useCase internal.UserUsecase
}

func NewUserEchoHandler(useCase internal.UserUsecase) *UserEchoHandler {
	return &UserEchoHandler{
		useCase: useCase,
	}
}

func (udh *UserEchoHandler) GetAllUsers(c echo.Context) error {
	users, status := udh.useCase.GetAllUsers()
	if status == models2.OK {
		return c.JSON(http.StatusOK, &users)
	}
	return pkg.NotFoundResponse(c)
}

func (udh *UserEchoHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	user := models2.User{ID: userID}
	log.Println("USER ME", userID.String())

	person := models2.User{}
	person.ID = user.ID

	s, errCode := udh.useCase.GetUser(person)
	if errCode != models2.OK {
		switch errCode {
		case models2.NotFound:
			return pkg.NotFoundResponse(c)
		default:
			return pkg.ServerErrorResponse(c, fmt.Errorf(""))
		}
	}

	return c.JSON(http.StatusOK, &s)
}

func (udh *UserEchoHandler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return pkg.BadRequestResponse(c, err)
	}

	person := models2.User{}
	person.ID = id

	foundPerson, errCode := udh.useCase.GetUser(person)
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

func (udh *UserEchoHandler) DeleteUser(c echo.Context) error {
	log.Println("DELETING")
	idStr := c.Param("id")
	v := validator.New()
	v.Check(idStr != "", "id", "must be provided")
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return pkg.BadRequestResponse(c, err)
	}

	person := models2.User{
		ID: id,
	}

	errCode := udh.useCase.DeleteUser(person)
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
	var input struct {
		Nickname string `json:"nickname,omitempty"`
		Name     string `json:"name,omitempty"`
		Surname  string `json:"surname,omitempty"`
		Password string `json:"-"`
	}
	err := c.Bind(&input)
	if err != nil {
		slog.Error(err.Error())
		return pkg.ServerErrorResponse(c, err)
	}

	person := &models2.User{
		Name:      input.Name,
		Nickname:  input.Nickname,
		Surname:   input.Surname,
		Activated: false,
	}

	err = person.Password.Set(input.Password)
	if err != nil {
		return pkg.ServerErrorResponse(c, err)
	}

	v := validator.New()
	validator.ValidateUserBySignUp(v, person)
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}

	s, errCode := udh.useCase.UpdateUser(*person)
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

func (udh *UserEchoHandler) FindUser(c echo.Context) error {
	name := c.QueryParam("name")
	v := validator.New()
	v.Check(name != "", "name", "must be provided")
	if !v.Valid() {
		return pkg.FailedValidationResponse(c, v.Errors)
	}
	slog.Info("Searching for " + name)

	users, errCode := udh.useCase.FindUser(name)
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

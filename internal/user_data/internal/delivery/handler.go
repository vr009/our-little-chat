package delivery

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	models2 "our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
	"our-little-chatik/internal/user_data/internal"
)

type UserdataEchoHandler struct {
	useCase internal.UserdataUseCase
}

func NewUserdataEchoHandler(useCase internal.UserdataUseCase) *UserdataEchoHandler {
	return &UserdataEchoHandler{
		useCase: useCase,
	}
}

func (udh *UserdataEchoHandler) GetAllUsers(c echo.Context) error {
	users, status := udh.useCase.GetAllUsers()
	if status == models2.OK {
		return c.JSON(http.StatusOK, &users)
	}
	return pkg.HandleErrorCode(status, models2.Error{Msg: "internal issue"}, c)
}

func (udh *UserdataEchoHandler) CreateUser(c echo.Context) error {
	person := models2.UserData{}
	err := c.Bind(&person)
	if err != nil {
		slog.Error(err.Error())
		return pkg.HandleErrorCode(models2.BadRequest, models2.Error{Msg: "bad body"}, c)
	}

	slog.Info("Unmarshalled:", "person", slog.AnyValue(person))

	newPerson, errCode := udh.useCase.CreateUser(person)
	if errCode != models2.OK {
		return pkg.HandleErrorCode(errCode, models2.Error{Msg: "Failed to create user"}, c)
	}
	slog.Info("Created: ", slog.AnyValue(newPerson))
	return c.JSON(http.StatusCreated, &newPerson)
}

func (udh *UserdataEchoHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	user := models2.User{UserID: userID}
	log.Println("USER ME", userID.String())

	person := models2.UserData{}
	person.UserID = user.UserID

	s, errCode := udh.useCase.GetUser(person)
	if errCode != models2.OK {
		return pkg.HandleErrorCode(errCode, models2.Error{Msg: "Failed to find info about the user"}, c)
	}

	return c.JSON(http.StatusOK, &s)
}

func (udh *UserdataEchoHandler) GetUser(c echo.Context) error {
	userID := c.Get("user_id").(uuid.UUID)
	user := models2.User{UserID: userID}

	idStr := c.QueryParam("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return pkg.HandleErrorCode(models2.BadRequest, models2.Error{Msg: "Bad id format"}, c)
	}

	slog.Info("requested info", "from user", user.UserID.String(), "about", idStr)

	person := models2.UserData{}
	person.UserID = id

	log.Println("SEARCHING FOR !!!!!!!!!!!!!!!!", id.String())
	foundPerson, errCode := udh.useCase.GetUser(person)
	if errCode != models2.OK {
		return pkg.HandleErrorCode(errCode, models2.Error{Msg: "Failed to get info about the user"}, c)
	}
	log.Println("FOUND !!!!!!!!!!!!!!!!", foundPerson.UserID.String())
	return c.JSON(http.StatusOK, &foundPerson)
}

func (udh *UserdataEchoHandler) DeleteUser(c echo.Context) error {
	log.Println("DELETING")
	person := models2.UserData{}
	err := c.Bind(&person)
	if err != nil {
		return pkg.HandleErrorCode(models2.BadRequest, models2.Error{Msg: "bad body"}, c)
	}

	errCode := udh.useCase.DeleteUser(person)
	if errCode != models2.OK {
		return pkg.HandleErrorCode(errCode, models2.Error{Msg: "Failed to delete user"}, c)
	}
	return c.JSON(http.StatusOK, models2.Error{Msg: "OK"})
}

func (udh *UserdataEchoHandler) UpdateUser(c echo.Context) error {
	person := models2.UserData{}
	err := c.Bind(&person)
	if err != nil {
		slog.Error(err.Error())
		return pkg.HandleErrorCode(models2.BadRequest, models2.Error{Msg: "bad body"}, c)
	}

	s, errCode := udh.useCase.UpdateUser(person)
	if errCode != models2.OK {
		return pkg.HandleErrorCode(errCode, models2.Error{Msg: "User update failed"}, c)
	}

	return c.JSON(http.StatusOK, &s)
}

func (udh *UserdataEchoHandler) CheckUserData(c echo.Context) error {
	person := models2.UserData{}
	err := c.Bind(&person)
	if err != nil {
		return pkg.HandleErrorCode(models2.BadRequest, models2.Error{Msg: "bad body"}, c)
	}

	_, errCode := udh.useCase.CheckUser(person)
	if errCode != models2.OK {
		return pkg.HandleErrorCode(errCode, models2.Error{Msg: "User data is inaccessible"}, c)
	}

	return c.JSON(http.StatusOK, models2.Error{Msg: "OK"})
}

func (udh *UserdataEchoHandler) FindUser(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return pkg.HandleErrorCode(models2.BadRequest, models2.Error{Msg: "bad parameter"}, c)
	}
	slog.Info("Searching for " + name)

	users, errCode := udh.useCase.FindUser(name)
	if errCode != models2.OK {
		return pkg.HandleErrorCode(errCode, models2.Error{Msg: "Inaccessible user data"}, c)
	}

	return c.JSON(http.StatusOK, users)
}

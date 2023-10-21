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

// GetMe godoc
// @Summary Get info about token holder.
// @Description get info about token holder.
// @Produce json
// @Tags users
// @Success 200 {object} models.HttpResponse
// @Failure 404 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /user/me [get]
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

	response := models2.EnvelopIntoHttpResponse(me, "me", http.StatusOK)
	return c.JSON(http.StatusOK, &response)
}

// GetUserForID godoc
// @Summary Get user for its id.
// @Description get user for its id.
// @Produce json
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} models.HttpResponse
// @Failure 404 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /user/{id} [get]
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
	response := models2.EnvelopIntoHttpResponse(foundPerson, "found_person", http.StatusOK)
	return c.JSON(http.StatusOK, &response)
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
	return c.JSON(http.StatusOK, &models2.HttpResponse{Message: "OK"})
}

// UpdateUser godoc
// @Summary Update user info.
// @Description update user info.
// @Accept json
// @Produce json
// @Tags users
// @Param request body models.UpdateUserRequest true "update user request"
// @Success 200 {object} models.HttpResponse
// @Failure 422 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /user/me [patch]
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

	response := models2.EnvelopIntoHttpResponse(s, "updated_person", http.StatusOK)
	return c.JSON(http.StatusOK, &response)
}

// SearchUsers godoc
// @Summary Search user for its nickname.
// @Description search user for its nickname.
// @Produce json
// @Tags users
// @Param offset query string true "nickname"
// @Success 200 {object} models.HttpResponse
// @Failure 404 {object} models.HttpResponse
// @Failure 500 {object} models.HttpResponse
// @Router /user/search [get]
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

	response := models2.EnvelopIntoHttpResponse(users, "found_users", http.StatusOK)
	return c.JSON(http.StatusOK, response)
}

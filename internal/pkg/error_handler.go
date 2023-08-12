package pkg

import (
	"github.com/labstack/echo/v4"
	"net/http"
	models2 "our-little-chatik/internal/models"
)

func HandleErrorCode(errCode models2.StatusCode, obj models2.Error, c echo.Context) error {
	switch errCode {
	case models2.NotFound:
		return c.JSON(http.StatusNotFound, &obj)
	case models2.InternalError:
		return c.JSON(http.StatusInternalServerError, &obj)
	case models2.BadRequest:
		return c.JSON(http.StatusBadRequest, &obj)
	case models2.Forbidden:
		return c.JSON(http.StatusForbidden, &obj)
	default:
		return c.JSON(http.StatusInternalServerError, &obj)
	}
}

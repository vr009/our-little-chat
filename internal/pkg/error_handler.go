package pkg

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"our-little-chatik/internal/models"
)

// The logError() method is a generic helper for logging an error message. Later in the // book we'll upgrade this to use structured logging, and record additional information // about the request including the HTTP method and URL.
func logError(err error) {
	log.Println(err)
}

// The ErrorResponse method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using an interface{}
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func ErrorResponse(c echo.Context, status int, message interface{}) error {
	env := models.EnvelopIntoHttpResponse(message, "description", status)
	// Write the response using the writeJSON() helper. If this happens to return an // error then log it, and fall back to sending the client an empty response with a // 500 Internal Server Error status code.
	err := c.JSON(status, &env)
	if err != nil {
		logError(c.NoContent(http.StatusInternalServerError))
	}
	return err
}

// The ServerErrorResponse method will be used when our application encounters an // unexpected problem at runtime. It logs the detailed error message, then uses the // ErrorResponse() helper to send a 500 Internal Server Error status code and JSON // response (containing a generic error message) to the client.
func ServerErrorResponse(c echo.Context, err error) error {
	logError(err)
	message := "the server encountered a problem and could not process your request"
	return ErrorResponse(c, http.StatusInternalServerError, message)
}

// The NotFoundResponse method will be used to send a 404 Not Found status code and // JSON response to the client.
func NotFoundResponse(c echo.Context) error {
	message := "the requested resource could not be found"
	return ErrorResponse(c, http.StatusNotFound, message)
}

// The MethodNotAllowedResponse method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func MethodNotAllowedResponse(c echo.Context) error {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request().Method)
	return ErrorResponse(c, http.StatusMethodNotAllowed, message)
}

// FailedValidationResponse Note that the errors parameter here has the type map[string]string,
// which is exactly the same as the errors map contained in our Validator type.
func FailedValidationResponse(c echo.Context, errors map[string]string) error {
	return ErrorResponse(c, http.StatusUnprocessableEntity, errors)
}

func BadRequestResponse(c echo.Context, err error) error {
	return ErrorResponse(c, http.StatusBadRequest, err.Error())
}

func InvalidCredentialsResponse(c echo.Context) error {
	message := "invalid authentication credentials"
	return ErrorResponse(c, http.StatusUnauthorized, message)
}

func UnauthorizedResponse(c echo.Context, err error) error {
	return ErrorResponse(c, http.StatusUnauthorized, err.Error())
}

func ForbiddenResponse(c echo.Context, err error) error {
	return ErrorResponse(c, http.StatusForbidden, err.Error())
}

package middleware

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func getAdminCreds() (user, pswd string) {
	user = os.Getenv("ADMIN_USER")
	pswd = os.Getenv("ADMIN_PASSWORD")
	return
}

// AdminAuth is the middleware function that .
func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	user, pswd := getAdminCreds()
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(user)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(pswd)) == 1 {
			return true, nil
		}
		return false, nil
	})(next)
}

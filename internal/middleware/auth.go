package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"our-little-chatik/internal/models"
	"our-little-chatik/internal/pkg"
)

// Auth is the middleware function that .
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken := c.Get("user").(*jwt.Token)
		claims := userToken.Claims.(*pkg.JwtCustomClaims)
		userIDStr := claims.UserID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusForbidden, models.Error{Msg: "bad token"})
		}

		c.Set("user_id", userID)
		return next(c)
	}
}

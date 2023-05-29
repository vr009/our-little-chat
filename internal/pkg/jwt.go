package pkg

import (
	"fmt"
	"os"
	"time"

	"our-little-chatik/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
)

func GetSignedKey() (string, error) {
	key := os.Getenv("JWT_SIGNED_KEY")
	if key == "" {
		return "", fmt.Errorf("no env provided")
	}
	return key, nil
}

func GenerateJWTToken(user models.User, expireCookie bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["UserID"] = user.UserID.String()
	if expireCookie {
		claims["exp"] = time.Now().Add(-time.Hour * 30).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Hour * 30).Unix()
	}

	mySignedKey, err := GetSignedKey()
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString([]byte(mySignedKey))
	if err != nil {
		slog.Error("generating JWT Token failed")
		return "", err
	}

	return tokenString, nil
}

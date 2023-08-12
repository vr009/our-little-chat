package pkg

import (
	"fmt"
	"os"
	"time"

	"our-little-chatik/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// JwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GetSignedKey() (string, error) {
	key := os.Getenv("JWT_SIGNED_KEY")
	if key == "" {
		return "", fmt.Errorf("no env provided")
	}
	return key, nil
}

func GenerateJWTTokenV2(user models.User, expireCookie bool) (string, error) {
	var m int
	if expireCookie {
		m = -1
	} else {
		m = 1
	}
	// Set custom claims
	claims := &JwtCustomClaims{
		user.UserID.String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72 * time.Duration(m))),
		},
	}
	mySignedKey, err := GetSignedKey()
	if err != nil {
		return "", err
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(mySignedKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

package pkg

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"our-little-chatik/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthHook(r *http.Request) (*models.User, error) {
	session := models.Session{}

	cookie, err := r.Cookie("Token")
	session.Token = cookie.Value
	if session.Token == "" {
		err = fmt.Errorf("no cookie provided")
		return nil, err
	}
	user := models.User{}

	token, err := jwt.Parse(session.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}

		mySigningKey, err := GetSignedKey()
		if err != nil {
			return "", err
		}

		return []byte(mySigningKey), nil
	})

	if !token.Valid {
		return nil, fmt.Errorf("invalid token\n")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expFl := claims["exp"].(float64)
		exp := time.Unix(int64(expFl), 0)
		if time.Until(exp) <= 0 {
			return nil, fmt.Errorf("expired token")
		}

		authed := claims["authorized"].(bool)
		if !authed {
			return nil, fmt.Errorf("unauthorized token")
		}

		id, err := uuid.Parse(claims["UserID"].(string))
		if err != nil {
			return nil, fmt.Errorf("invalid token: %v\n", err)
		}
		user.UserID = id
		log.Println("PARSED", id)
	} else {
		return nil, fmt.Errorf("spoiled token\n")
	}

	return &user, nil
}
